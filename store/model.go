package lurkstore

import (
	"database/sql"
)

var db *sql.DB

// TODO make backend choice pluggable

func Init(filename string) (err error) {
	db, err = DbInit(filename)
	return
}

// Store a new registration
// Returns the unique id for the newly created record
func NewRegistration(csr string, lifetime uint) (string, error) {
	return DbAddRegistration(db, csr, lifetime)
}

// Return the Registration record associated to the supplied id, if found
func GetRegistrationById(id string) (*Registration, error) {
	return DbGetRegistrationById(db, id)
}

// Fetch the oldest registration in state "new" (if one exists) and mark it
// as "work-in-progress"
func DequeueRegistration() (*Registration, error) {
	return DbDequeueRegistration(db)
}

// Mark a work-in-progress as successfully completed
func UpdateSuccessfulRegistration(id string, certURL string, lifetime uint,
	ttl string) error {
	return DbUpdateSuccessfulRegistration(db, id, certURL, lifetime, ttl)
}

// Mark a work-in-progress as failed
// TODO control how long the registration is visible
func UpdateFailedRegistration(id string, errmsg string) error {
	return DbUpdateFailedRegistration(db, id, errmsg)
}

// Not part of the API -- diagnostics/introspection only
func ListRegistrations() ([]Registration, error) {
	return DbListRegistrations(db)
}
