package lurkstore

import (
	"database/sql"
	"errors"
)

var db *sql.DB

// TODO make backend choice pluggable
func Init(filename string) error {
	var err error

	db, err = DbInit(filename)
	if err != nil {
		return err
	}

	err = DbCreateRegistrationTable(db)
	if err != nil {
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

// Fetch the oldest registration in state "new"
func GetNewRegistration() (Registration, error) {
	// TODO
	return Registration{}, nil
}

// Finalise a work-in-progress
func FinaliseRegistration(id string, status string, certURL string, lifetime uint) error {
	// TODO
	return errors.New("hehe")
}
