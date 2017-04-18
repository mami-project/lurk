package lurkstore

import (
	"os"
	"testing"
)

const dbfile = "./test.db"

func TestNewRegistration(t *testing.T) {
	err := Init(dbfile)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}

	defer os.Remove(dbfile)

	_, err = NewRegistration("test csr", 1234)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}
}

func TestGetRegistrationById(t *testing.T) {
	_ = Init(dbfile)
	defer os.Remove(dbfile)

	id, err := NewRegistration("another csr", 7890)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}

	want := Registration{Status: "new", CSR: "another csr", Lifetime: 7890}

	got, err := GetRegistrationById(id)
	if err != nil {
		t.Errorf("GetRegistrationById returned %v", err)
	}

	if got.CSR != want.CSR {
		t.Errorf("CSR mismatch: got %s, want %s", got.CSR, want.CSR)
	}

	if got.Lifetime != want.Lifetime {
		t.Errorf("Lifetime mismatch: got %d, want %d", got.Lifetime, want.Lifetime)
	}
}

func TestWorkOnNewRegistration(t *testing.T) {
	_ = Init(dbfile)
	defer os.Remove(dbfile)

	wanted_csr := "an older csr not yet processed"

	_, err := NewRegistration(wanted_csr, 1010)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}

	_, err = NewRegistration("a newer csr not yet processed", 2020)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}

	got, err := WorkOnNewRegistration()
	if err != nil {
		t.Errorf("WorkOnNewRegistration returned %v", err)
	}

	if got.CSR != wanted_csr {
		t.Errorf("CSR mismatch: got %s, want %s)", got.CSR, wanted_csr)
	}

	// TODO check status is now "work-in-progress"
}

func TestUpdateSuccessfulRegistration(t *testing.T) {
	_ = Init(dbfile)
	defer os.Remove(dbfile)

	// TODO
}
