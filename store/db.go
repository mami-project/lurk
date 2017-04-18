package lurkstore

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func DbInit(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	return db, nil
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
		certURL   TEXT NOT NULL DEFAULT ""

		CHECK (status IN ("new", "wip", "done", "failed")),
		CHECK (lifetime > 0)
	);
	`

	_, err := db.Exec(sql_query)

	return err
}

// Return the (unique) identifier associated to the added record, or the empty
// string on error
func DbAddRegistration(db *sql.DB, csr string, lifetime uint) (string, error) {
	sql_query := "INSERT INTO registration(csr, lifetime) VALUES(?, ?)"

	stmt, err := db.Prepare(sql_query)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	res, err := stmt.Exec(csr, lifetime)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func DbGetRegistrationById(db *sql.DB, id string) (*Registration, error) {
	sql_query := `
	SELECT id,
	       status,
	       csr,
	       created,
	       completed,
	       expires,
	       lifetime,
	       certURL
	  FROM registration
	 WHERE id = ?
	`

	reg := Registration{}

	err := db.QueryRow(sql_query, id).
		Scan(&reg.Id, &reg.Status, &reg.CSR, &reg.CreationDate,
			&reg.CompletionDate, &reg.ExpirationDate, &reg.Lifetime,
			&reg.CertURL)

	if err != nil {
		return nil, err
	}

	return &reg, nil
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
			 certURL
	    FROM registration
	   WHERE status = "new"
	ORDER BY created DESC
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
		&reg.ExpirationDate, &reg.Lifetime, &reg.CertURL)

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

// TODO add tests
func DbUpdateSuccessfulRegistration(db *sql.DB, id string, certURL string,
	lifetime uint) error {
	sql_query := `
	UPDATE registration
	   SET status = "done",
	       certURL = ?,
	       lifetime = ?,
		   completed = CURRENT_TIMESTAMP
	 WHERE id = ? AND
	       status = "wip"
	`

	stmt, err := db.Prepare(sql_query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(certURL, lifetime, id)
	if err != nil {
		return err
	}

	return nil
}

// TODO add tests
func DbUpdateFailedRegistration(db *sql.DB, id string) error {
	sql_query := `
	UPDATE registration
	   SET status = "failed",
		   completed = CURRENT_TIMESTAMP
	 WHERE id = ? AND
	       status = "wip"
	`

	stmt, err := db.Prepare(sql_query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
