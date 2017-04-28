package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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
	ct := "application/json"

	json := `{
		"csr": "...",
		"lifetime": 365
	}`
	body := []byte(json)

	req, _ := http.NewRequest("POST", NewRegistrationPath, bytes.NewBuffer(body))
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", ct)

	res := runTransaction(NewRouter(), req)

	// 201
	if res.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusCreated, res.Code)
	}

	// Location: https://star-proxy.example.net/star/registration/567
	loc := res.Header().Get("Location")
	if loc == "" {
		t.Errorf("no Location header")
	}

	wanted := "https://" + host + NewRegistrationPath + "/1"
	if loc != wanted {
		t.Errorf("want Location %s, got %s", wanted, loc)
	}
}

func TestFailNewRegistrationNoCSR(t *testing.T) {
	host := "star-proxy.example.net"
	ct := "application/json"

	json := `{
		"lifetime": 365
	}`
	body := []byte(json)

	req, _ := http.NewRequest("POST", NewRegistrationPath, bytes.NewBuffer(body))
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", ct)

	res := runTransaction(NewRouter(), req)

	t.Logf("%v", req)
	t.Logf("%v", res)

	// 400
	if res.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, res.Code)
	}
}
