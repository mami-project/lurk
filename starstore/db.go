package starstore

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// return a handle to an "open" DB
// the supplied dataSource is a driver-specific data source name;
// in case of SQLite, it's an existing local file.
func DbInit(dataSource string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		return
	}

	err = DbCreateRegistrationTable(db)

	return
}

func DbCreateRegistrationTable(db *sql.DB) error {
	sql_query := `
	CREATE TABLE IF NOT EXISTS registration (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		status    TEXT NOT NULL DEFAULT "new",
		csr       BLOB,
		created   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed DATETIME DEFAULT NULL,
		expires   DATETIME DEFAULT NULL,
		lifetime  INTEGER,
		certURL   TEXT NOT NULL DEFAULT "",
		errmsg    TEXT

		CHECK (status IN ("new", "wip", "done", "failed")),
		CHECK (lifetime > 0)
	)
	`

	_, err := db.Exec(sql_query)

	return err
}

// Return the (unique) identifier associated to the added record, or the empty
// string on error
func DbAddRegistration(db *sql.DB, r Registration) (string, error) {
	sql_query := "INSERT INTO registration(csr, lifetime) VALUES(?, ?)"

	stmt, err := db.Prepare(sql_query)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	res, err := stmt.Exec(r.CSR, r.Lifetime)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

// Return the (unique) identifier associated to the added record, or the empty
// string on error
func DbGetRegistrationById(db *sql.DB, id string) (*Registration, error) {
	sql_query := `
	SELECT id,
	       status,
	       csr,
	       created,
	       completed,
	       expires,
	       lifetime,
	       certURL,
		   errmsg
	  FROM registration
	 WHERE id = ?
	`

	var r Registration

	err := db.QueryRow(sql_query, id).
		Scan(&r.Id, &r.Status, &r.CSR, &r.CreationDate,
			&r.CompletionDate, &r.ExpirationDate, &r.Lifetime,
			&r.CertURL, &r.ErrMsg)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func DbDequeueRegistration(db *sql.DB) (reg *Registration, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sql_select_oldest_waiting := `
	  SELECT id,
	         status,
	         csr,
			 created,
			 completed,
			 expires,
			 lifetime,
			 certURL,
			 errmsg
	    FROM registration
	   WHERE status = "new"
	ORDER BY created ASC
	   LIMIT 1
	`

	// Get the oldest waiting registration; the query will return
	// at most one result.
	var rows *sql.Rows

	if rows, err = tx.Query(sql_select_oldest_waiting); err != nil {
		return
	}

	defer rows.Close()

	if !rows.Next() {
		if rows.Err() != nil {
			return
		}
		// Nothing to dequeue
		// TODO make sure reg is nil
		return
	}

	reg = new(Registration)

	// Copy results in to the Registration struct
	err = rows.Scan(&reg.Id, &reg.Status, &reg.CSR,
		&reg.CreationDate, &reg.CompletionDate,
		&reg.ExpirationDate, &reg.Lifetime, &reg.CertURL,
		&reg.ErrMsg)

	// Atomically update the status to work-in-progress
	sql_update_status_to_wip := `
	UPDATE registration
	   SET status = "wip"
	 WHERE id = ?
	`

	_, err = tx.Exec(sql_update_status_to_wip, reg.Id)
	if err != nil {
		return
	}

	// update reg.Status and return
	reg.Status = "wip"

	return
}

func DbUpdateSuccessfulRegistration(db *sql.DB, id string, certURL string,
	lifetime uint, ttl string) error {
	sql_query := `
	UPDATE registration
	   SET status = "done",
	       certURL = ?,
	       lifetime = ?,
		   completed = CURRENT_TIMESTAMP,
		   expires = datetime(CURRENT_TIMESTAMP, ?)
	 WHERE id = ? AND
	       status = "wip"
	`

	stmt, err := db.Prepare(sql_query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(certURL, lifetime, ttl, id)
	if err != nil {
		return err
	}

	return nil
}

func DbUpdateFailedRegistration(db *sql.DB, id string, errmsg string) error {
	sql_query := `
	UPDATE registration
	   SET status = "failed",
		   completed = CURRENT_TIMESTAMP,
		   errmsg = ?
	 WHERE id = ?
	`

	stmt, err := db.Prepare(sql_query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(errmsg, id)
	if err != nil {
		return err
	}

	return nil
}

func DbListRegistrations(db *sql.DB) ([]Registration, error) {
	// Include registrations that have not been processed yet
	// (i.e., "expires IS NULL")
	sql_query := `
	SELECT *
	  FROM registration
	 WHERE expires IS NULL OR
	       expires > CURRENT_TIMESTAMP;
	`

	rows, err := db.Query(sql_query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	registrations := []Registration{}

	for rows.Next() {
		var r Registration
		err := rows.Scan(&r.Id, &r.Status, &r.CSR, &r.CreationDate,
			&r.CompletionDate, &r.ExpirationDate, &r.Lifetime, &r.CertURL,
			&r.ErrMsg)
		if err != nil {
			return nil, err
		}
		registrations = append(registrations, r)
	}

	return registrations, nil
}

// Remove all rows and reset the auto-increment id
func DbRemoveAll(db *sql.DB) error {
	sql_query := `
	DELETE
	  FROM registration;
	DELETE
	  FROM sqlite_sequence
	 WHERE name='registration'
	`
	_, err := db.Exec(sql_query)

	return err
}
