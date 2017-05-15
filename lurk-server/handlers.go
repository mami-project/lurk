package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"../starstore"
)

const HelloSTAR string = "Hello STAR!"

// TODO move these to configuration
var DefaultSTARHost string = "todo-setme.example.net"
var PollIntervalInSeconds uint64 = 10

func internalError(res http.ResponseWriter, err string) {
	log.Printf("[ERROR]: %s", err)
	res.WriteHeader(http.StatusInternalServerError)
}

func replyError(res http.ResponseWriter, code int, message string) {
	body, err := json.Marshal(MsgError{message})
	if err != nil {
		internalError(res, fmt.Sprintf("json.Marshal failed: %s", err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(body)
}

// Expires header
func replyDone(res http.ResponseWriter, r *starstore.Registration) {
	body, err := json.Marshal(MsgRegistrationDone{"success", r.Lifetime, r.CertURL})
	if err != nil {
		internalError(res, fmt.Sprintf("json.Marshal failed: %s", err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

func replyFailed(res http.ResponseWriter, r *starstore.Registration) {
	details := ""
	if r.ErrMsg.Valid {
		details = r.ErrMsg.String
	}

	body, err := json.Marshal(MsgRegistrationFailed{"failed", details})
	if err != nil {
		internalError(res, fmt.Sprintf("json.Marshal failed: %s", err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

func maxAgeForPolling(pollInterval uint64) string {
	return "max-age=" + strconv.FormatUint(pollInterval-1, 10)
}

func replyPending(res http.ResponseWriter) {

	body, err := json.Marshal(MsgRegistrationPending{"pending", PollIntervalInSeconds})
	if err != nil {
		internalError(res, fmt.Sprintf("json.Marshal failed: %s", err))
		return
	}

	// Set an explicit C-C=max-age to make sure cache expiration heuristics
	// do not interfere with polling.
	res.Header().Set("Cache-Control", maxAgeForPolling(PollIntervalInSeconds))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

func assembleRegistrationURL(req *http.Request, id string) string {
	host := req.Host
	if host == "" {
		host = req.Header.Get("Host")
		if host == "" {
			host = DefaultSTARHost
		}
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

	id, err := starstore.NewRegistration(r)
	if err != nil {
		replyError(res, http.StatusBadRequest, err.Error())
	}

	res.Header().Set("Location", assembleRegistrationURL(req, id))
	res.WriteHeader(http.StatusCreated)
}

// Create a new registration object

func PollRegistrationStatus(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	r, err := starstore.GetRegistrationById(vars["id"])
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
		replyFailed(res, r)
	}
}

// Return the list of all registration requests (debug-only)
func RegistrationsList(res http.ResponseWriter, req *http.Request) {
	var rs []starstore.Registration

	rs, err := starstore.ListRegistrations()
	if err != nil {
		internalError(res, fmt.Sprintf("list registration failed: %s", err))
		return
	}

	body, err := json.Marshal(rs)
	if err != nil {
		internalError(res, fmt.Sprintf("json.Marshal failed: %s", err))
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}
