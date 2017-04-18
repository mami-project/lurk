package lurkstore

import (
	"database/sql"
	"os"
)

var db *sql.DB

// TODO make backend choice pluggable
func Init(filename string) error {
	var err error

	db, err = DbInit(filename)
	if err != nil {
		_ = os.Remove(filename)
		return err
	}

	err = DbCreateRegistrationTable(db)
	if err != nil {
		_ = os.Remove(filename)
		return err
	}

	return nil
}

// Store new registration
// Returns the unique id for the newly created record
func NewRegistration(csr string, lifetime uint) (string, error) {
	return DbAddRegistration(db, csr, lifetime)
}

// Return the Registration record associated to the supplied id, if found
func GetRegistrationById(id string) (*Registration, error) {
	return DbGetRegistrationById(db, id)
}

// Fetch the oldest registration in state "new" and mark it as "work-in-progress"
func WorkOnNewRegistration() (*Registration, error) {
	return DbGetNewRegistration(db)
}

// Mark a work-in-progress as successfully completed
// TODO ttl of the registration
func UpdateSuccessfulRegistration(id string, certURL string, lifetime uint) error {
	return DbUpdateSuccessfulRegistration(db, id, certURL, lifetime)
}

// Mark a work-in-progress as failed
// TODO ttl of the registration
func UpdateFailedRegistration(id string) error {
	return DbUpdateFailedRegistration(db, id)
}
