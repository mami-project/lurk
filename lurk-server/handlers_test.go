package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"../starstore"

	"github.com/gorilla/mux"
)

// handle setup/teardown (we need to initialise the store)
func runTests(m *testing.M) int {
	dbfile := "./test.db"

	err := starstore.Init(dbfile)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	defer os.Remove(dbfile)

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTransaction(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func postRegistration(host string) *httptest.ResponseRecorder {
	ct := "application/json"

	json := `{
		"csr": "...",
		"lifetime": 365
	}`
	body := []byte(json)

	req, _ := http.NewRequest("POST", NewRegistrationPath, bytes.NewBuffer(body))
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", ct)

	return runTransaction(NewRouter(), req)
}

func TestIndex(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := runTransaction(NewRouter(), req)

	// 200
	if res.Code != http.StatusOK {
		t.Errorf("want 200, got %d", res.Code)
	}

	want := HelloSTAR
	got, _ := ioutil.ReadAll(res.Body)

	// "Hello STAR!"
	if string(got) != want {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestCreateNewRegistration(t *testing.T) {
	starstore.RemoveAllRegistrations()

	host := "star-proxy.example.net"

	res := postRegistration(host)

	// 201
	if res.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusCreated, res.Code)
	}

	// Location: https://star-proxy.example.net/star/registration/567
	loc := res.Header().Get("Location")
	if loc == "" {
		t.Errorf("no Location header")
	}

	wanted := "http://" + host + NewRegistrationPath + "/1"
	if loc != wanted {
		t.Errorf("want Location %s, got %s", wanted, loc)
	}
}

func TestFailNewRegistrationNoCSR(t *testing.T) {
	host := "star-proxy.example.net"
	ct := "application/json"

	j := `{
		"lifetime": 365
	}`
	body := []byte(j)

	req, _ := http.NewRequest("POST", NewRegistrationPath, bytes.NewBuffer(body))
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", ct)

	res := runTransaction(NewRouter(), req)

	// 400
	if res.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, res.Code)
	}

	// json.error == "empty CSR"
	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["error"] != "empty CSR" {
		t.Errorf("want error=empty CSR, got error=%s", m["error"])
	}
}

func TestPollRegistrationStatusIdNotFound(t *testing.T) {
	starstore.RemoveAllRegistrations()

	host := "star-proxy.example.net"

	req, _ := http.NewRequest("GET", "/star/registration/123456", nil)
	req.Header.Set("Host", host)

	res := runTransaction(NewRouter(), req)

	// 400
	if res.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, res.Code)
	}

	// json.error == "sql: no rows in result set"
	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["error"] != "sql: no rows in result set" {
		t.Errorf("want error=sql: no rows in result set, got error=%s", m["error"])
	}
}

func TestPollRegistrationStatusPending(t *testing.T) {
	starstore.RemoveAllRegistrations()

	res := postRegistration("star-proxy.example.net")

	// 201
	if res.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusOK, res.Code)
	}

	loc := res.Header().Get("Location")
	if loc == "" {
		t.Errorf("no Location header")
	}

	req, _ := http.NewRequest("GET", loc, nil)
	res = runTransaction(NewRouter(), req)

	// 200
	if res.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, res.Code)
	}

	// Retry-After
	ra := res.Header().Get("Retry-After")
	if ra == "" {
		t.Errorf("no Retry-After header")
	}

	// json.status == "pending"
	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["status"] != "pending" {
		t.Errorf("want status=pending, got status=%s", m["status"])
	}
}

func TestPollRegistrationStatusFailed(t *testing.T) {
	starstore.RemoveAllRegistrations()

	// Create new registration
	res := postRegistration("star-proxy.example.net")

	// 201
	if res.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusOK, res.Code)
	}

	loc := res.Header().Get("Location")
	if loc == "" {
		t.Errorf("no Location header")
	}

	// Extract the id of the registration from the returned location URL
	s := strings.Split(loc, "/")
	id := s[len(s)-1]

	// Fail the registration
	err := starstore.UpdateFailedRegistration(id, "a detailed description of the failure")
	if err != nil {
		t.Errorf("unable to fail the transaction: %s", err)
	}

	// Now retrieve the registration and verify its status is failed
	req, _ := http.NewRequest("GET", loc, nil)
	res = runTransaction(NewRouter(), req)

	// 200
	if res.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, res.Code)
	}

	// json.status == "failed"
	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["status"] != "failed" {
		t.Errorf("want status=failed, got status=%s", m["status"])
	}
}

// TODO cleanup -
func TestPollRegistrationStatusSuccess(t *testing.T) {
	starstore.RemoveAllRegistrations()

	// Create new registration
	res := postRegistration("star-proxy.example.net")

	// 201
	if res.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusOK, res.Code)
	}

	loc := res.Header().Get("Location")
	if loc == "" {
		t.Errorf("no Location header")
	}

	// Extract the id of the registration from the returned location URL
	s := strings.Split(loc, "/")
	id := s[len(s)-1]

	// Dequeue the registration (internally this moves its status from new to wip)
	r, err := starstore.DequeueRegistration()
	if err != nil {
		t.Errorf("could not dequeue registration: %s", err)
	}

	if r.Id != id {
		t.Errorf("want id %d, got id %d", id, r.Id)
	}

	// Mark the registration as successfully completed
	certURL := "http://acme.example.com/a-cert"
	err = starstore.UpdateSuccessfulRegistration(id, certURL, 366, "+3 days")
	if err != nil {
		t.Errorf("unable to complete the transaction: %s", err)
	}

	// Now retrieve the registration and verify it has a successful status
	req, _ := http.NewRequest("GET", loc, nil)
	res = runTransaction(NewRouter(), req)

	// 200
	if res.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, res.Code)
	}

	// json.status == "success"
	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["status"] != "success" {
		t.Errorf("want status success, got %s", m["status"])
	}
	// XXX lifetime is a number, not a string
	if m["lifetime"] != "366" {
		t.Errorf("want lifetime 366, got %s", m["lifetime"])
	}
	if m["certificates"] != certURL {
		t.Errorf("want certificates %s, got %s", certURL, m["certificates"])
	}
}
