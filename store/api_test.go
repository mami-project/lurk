package lurkstore

import (
	"testing"
)

func TestNewRegistration(t *testing.T) {
	err := Init("test")
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}

	_, err = NewRegistration("test csr", 1234)
	if err != nil {
		t.Errorf("NewRegistration returned %v", err)
	}
}

func TestGetRegistrationById(t *testing.T) {
	Init("test")

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
		t.Errorf("CSR mismatch")
	}

	if got.Lifetime != want.Lifetime {
		t.Errorf("CSR mismatch")
	}
}

// TODO
func TestGetNewRegistration(t *testing.T) {
	Init("test")

	want := Registration{}
	got, _ := GetNewRegistration()
	if got != want {
		t.Errorf("Got: %r, want: %r.", got, want)
	}
}

// TODO
func TestFinaliseRegistration(t *testing.T) {
	Init("test")

	err := FinaliseRegistration("an id", "done", "http://acme.example.com/wxyz/crt", 1234)
	if err.Error() == "eheh" {
		t.Errorf("Expecting error")
	}
}
