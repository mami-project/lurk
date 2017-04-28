package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../starstore"
)

const HelloSTAR string = "Hello STAR!"

// TODO configuration
var DefaultSTARHost string = "todo-setme.example.net"

func replyError(res http.ResponseWriter, code int, message string) {
	body, _ := json.Marshal(map[string]string{"error": message})

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(body)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, HelloSTAR)
}

func registrationURL(req *http.Request, id string) string {
	host := req.Header.Get("Host")
	if host == "" {
		host = DefaultSTARHost
	}

	return "https://" + host + req.URL.Path + "/" + id
}

func CreateNewRegistration(res http.ResponseWriter, req *http.Request) {
	var r starstore.Registration

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		replyError(res, http.StatusBadRequest, err.Error())
	}

	defer req.Body.Close()

	id, err := r.NewRegistration()
	if err != nil {
		replyError(res, http.StatusBadRequest, err.Error())
	}

	res.Header().Set("Location", registrationURL(req, id))
	res.WriteHeader(http.StatusCreated)
}

// Return the list of all registration requests
func RegistrationsList(w http.ResponseWriter, r *http.Request) {
	// TODO(tho)
}

// Create a new registration object

func RegistrationProgress(w http.ResponseWriter, r *http.Request) {
	// TODO(tho)
}
