package lurkstore

import (
	"testing"
)

func TestNewRegistration(t *testing.T) {
	id, _ := NewRegistration("test csr", 1234)
	if id != "id" {
		t.Errorf("Got: %s, want: id.", id)
	}
}

func TestGetRegistrationById(t *testing.T) {
	want := Registration{}
	got, _ := GetRegistrationById("an id")
	if got != want {
		t.Errorf("Got: %r, want: %r.", got, want)
	}
}

func TestGetNewRegistration(t *testing.T) {
	want := Registration{}
	got, _ := GetNewRegistration()
	if got != want {
		t.Errorf("Got: %r, want: %r.", got, want)
	}
}

func TestFinaliseRegistration(t *testing.T) {
	err := FinaliseRegistration("an id", "done", "http://acme.example.com/wxyz/crt", 1234)
	if err.Error() == "eheh" {
		t.Errorf("Expecting error")
	}
}
