package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

const (
	NewRegistrationPath   = "/star/registration"
	ListRegistrationsPath = "/star/registrations"
)

var routes = []Route{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: Index,
	},
	Route{
		Name:        "Index of fresh STAR Registrations",
		Method:      "GET",
		Pattern:     ListRegistrationsPath,
		HandlerFunc: RegistrationsList,
	},
	Route{
		Name:        "Poll status of a STAR Registration",
		Method:      "GET",
		Pattern:     "/star/registration/{id}",
		HandlerFunc: PollRegistrationStatus,
	},
	Route{
		Name:        "Request new STAR Registration",
		Method:      "POST",
		Pattern:     NewRegistrationPath,
		HandlerFunc: CreateNewRegistration,
	},
}
