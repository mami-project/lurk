package lurkstore

import "errors"

// Store new registration
// Returns the unique id for the newly created record
func NewRegistration(csr string, lifetime uint) (string, error) {
	// TODO
	return "id", nil
}

// Return the Registration record associated to the supplied id, if found
func GetRegistrationById(id string) (Registration, error) {
	// TODO
	return Registration{}, nil
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
