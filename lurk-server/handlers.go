package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"../starstore"
)

const HelloSTAR string = "Hello STAR!"

// TODO configuration
var DefaultSTARHost string = "todo-setme.example.net"
var PollIntervalInSeconds string = "10"

func replyError(res http.ResponseWriter, code int, message string) {
	m := map[string]string{
		"error": message,
	}

	body, _ := json.Marshal(m)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(body)
}

// TODO lifetime is a number, not a string
// Expires header
func replyDone(res http.ResponseWriter, r starstore.Registration) {
	m := map[string]string{
		"status":       "success",
		"lifetime":     strconv.FormatUint(uint64(r.Lifetime), 10),
		"certificates": r.CertURL,
	}

	body, _ := json.Marshal(m)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

// TODO add failure details?
func replyFailed(res http.ResponseWriter) {
	m := map[string]string{
		"status": "failed",
	}

	body, _ := json.Marshal(m)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

// TODO add c-c max-age (depends on retry-after)
func replyPending(res http.ResponseWriter) {
	m := map[string]string{
		"status": "pending",
	}

	body, _ := json.Marshal(m)

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Retry-After", PollIntervalInSeconds)
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

func registrationURL(req *http.Request, id string) string {
	// XXX this seems to contradict documentation:
	// "For incoming requests, the Host header is promoted to the
	//  Request.Host field and removed from the Header map."
	host := req.Header.Get("Host")
	if host == "" {
		host = DefaultSTARHost
	}

	scheme := req.URL.Scheme
	if scheme == "" {
		if req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	return scheme + "://" + host + req.URL.Path + "/" + id
}

// GET /
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, HelloSTAR)
}

// POST /star/registration
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

// Create a new registration object

func PollRegistrationStatus(res http.ResponseWriter, req *http.Request) {
	var r starstore.Registration

	vars := mux.Vars(req)

	err := r.GetRegistrationById(vars["id"])
	if err != nil {
		replyError(res, http.StatusBadRequest, err.Error())
	}

	switch r.Status {
	case "new":
		fallthrough
	case "wip":
		replyPending(res)
	case "done":
		replyDone(res, r)
	case "failed":
		replyFailed(res)
	}
}

// Return the list of all registration requests
func RegistrationsList(w http.ResponseWriter, r *http.Request) {
	// TODO(tho)
}
