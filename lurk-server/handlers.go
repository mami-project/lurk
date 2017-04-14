package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello LURK!\n")
}

// Return the list of all registration requests
func RegistrationsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(regs); err != nil {
		panic(err)
	}
}

// Create a new registration object
func RegistrationCreate(w http.ResponseWriter, r *http.Request) {
	var reg LurkRegistration

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &reg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := CreateLurkRegistration(reg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// done
	w.WriteHeader(http.StatusCreated)
}

func RegistrationProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var registrationId int
	var err error
	if registrationId, err = strconv.Atoi(vars["registrationId"]); err != nil {
		// TODO(tho) better error handling?
		panic(err)
	}

	reg := LookupLurkRegistration(registrationId)

	if reg.Id > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(reg); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}
