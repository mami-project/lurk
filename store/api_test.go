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

func TestDequeueRegistration(t *testing.T) {
	_ = Init(dbfile)
	defer os.Remove(dbfile)

	wanted_csr1 := "an older csr not yet processed"
	wanted_csr2 := "a newer csr not yet processed"
	wanted_status := "wip"

	_, err := NewRegistration(wanted_csr1, 1010)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}

	_, err = NewRegistration(wanted_csr2, 2020)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}

	// Dequeue the first registration request
	got, err := DequeueRegistration()
	if err != nil {
		t.Errorf("DequeueRegistration returned %v", err)
	}

	if got.CSR != wanted_csr1 {
		t.Errorf("CSR mismatch: got %s, want %s)", got.CSR, wanted_csr1)
	}

	if got.Status != wanted_status {
		t.Errorf("Status mismatch: got %s, want %s)", got.Status, wanted_status)
	}

	// Dequeue the second registration request
	got, err = DequeueRegistration()
	if err != nil {
		t.Errorf("DequeueRegistration returned %v", err)
	}

	if got.CSR != wanted_csr2 {
		t.Errorf("CSR mismatch: got %s, want %s)", got.CSR, wanted_csr2)
	}

	if got.Status != wanted_status {
		t.Errorf("Status mismatch: got %s, want %s)", got.Status, wanted_status)
	}

	got, err = DequeueRegistration()
	if got != nil || err != nil {
		t.Errorf("Expecting (nil, nil), got (%v, %v)", got, err)
	}
}

//func TestUpdateSuccessfulRegistration(t *testing.T) {
//	_ = Init(dbfile)
//	defer os.Remove(dbfile)
//
//	// TODO
//}
