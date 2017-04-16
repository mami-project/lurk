package lurkstore

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func DbInit(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	return db
}

func DbCreateRegistrationTable(db *sql.DB) error {
	sql_table := `
	CREATE TABLE IF NOT EXISTS registration (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		status    TEXT NOT NULL DEFAULT "new",
		csr       BLOB,
		created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed TIMESTAMP DEFAULT NULL,
		expires   TIMESTAMP DEFAULT NULL,
		lifetime  INTEGER,
		certURL   TEXT NOT NULL DEFAULT ""

		CHECK (status IN ("new", "wip", "done", "failed")),
		CHECK (lifetime > 0)
	);
	`

	_, err := db.Exec(sql_table)

	return err
}

// Return the (unique) identifier associated to the added record
func DbAddRegistration(db *sql.DB, csr string, lifetime uint) (string, error) {
	sql_additem := "INSERT INTO registration(csr, lifetime) VALUES(?, ?)"

	stmt, err := db.Prepare(sql_additem)
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
	sql_get_registration_by_id := "SELECT status, certURL FROM registration WHERE Id = ?"

	reg := Registration{}

	err := db.QueryRow(sql_get_registration_by_id, id).
		Scan(&reg.Status, &reg.CertURL)

	if err != nil {
		return nil, err
	}

	return &reg, nil
}
