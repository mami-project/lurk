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
		"csr": "5jNudRx6Ye4HzKEqT5...FS6aKdZeGsysoCo4H9P",
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

	var m MsgRegistrationPending
	err := json.Unmarshal(res.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("decoding MsgRegistrationPending failed: %s", err)
	}

	// json.status == "pending"
	if m.Status != "pending" {
		t.Errorf("want status=pending, got status=%s", m.Status)
	}

	if m.RetryAfter != PollIntervalInSeconds {
		t.Errorf("want retry-after=%d, got %s", PollIntervalInSeconds, m.Status)
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
	details := "a detailed description of the failure"
	err := starstore.UpdateFailedRegistration(id, details)
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

	var m MsgRegistrationFailed
	err = json.Unmarshal(res.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("decoding MsgRegistrationFailed failed: %s", err)
	}

	// json.status == "failed"
	if m.Status != "failed" {
		t.Errorf("want status=failed, got status=%s", m.Status)
	}

	if m.Details != details {
		t.Errorf("want details %s, got %s", m.Details)
	}
}

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
		t.Errorf("dequeuing: want id %d, got id %d", id, r.Id)
	}

	// Mark the registration as successfully completed
	certURL := "http://acme.example.com/a-cert"
	lifetime := uint(366)
	err = starstore.UpdateSuccessfulRegistration(id, certURL, lifetime, "+3 days")
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

	var m MsgRegistrationDone
	err = json.Unmarshal(res.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("decoding MsgRegistrationDone failed: %s", err)
	}

	// json.status == "success"
	if m.Status != "success" {
		t.Errorf("want status success, got %s", m.Status)
	}

	if m.Lifetime != lifetime {
		t.Errorf("want lifetime %d, got %d", lifetime, m.Lifetime)
	}

	if m.CertURL != certURL {
		t.Errorf("want certificates %s, got %s", certURL, m.CertURL)
	}
}
