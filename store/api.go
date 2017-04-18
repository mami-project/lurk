package lurkstore

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

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

// TODO make backend choice pluggable
func Init2(filename string) (err error) {
	db, err := DbInit(filename)
	if err != nil {
		return
	}

	fmt.Println("DB: %v", db)

	defer func() {
		if err != nil {
			_ = os.Remove(filename)
			return
		}
		return
	}()

	err = DbCreateRegistrationTable(db)

	fmt.Println("ERR: %v", err)

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
// TODO ttl of the registration
func UpdateSuccessfulRegistration(id string, certURL string, lifetime uint) error {
	return DbUpdateSuccessfulRegistration(db, id, certURL, lifetime)
}

// Mark a work-in-progress as failed
// TODO ttl of the registration
func UpdateFailedRegistration(id string) error {
	return DbUpdateFailedRegistration(db, id)
}
