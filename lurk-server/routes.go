package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

const NewRegistrationPath = "/star/registration"

var routes = []Route{
	Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: Index,
	},
	Route{
		Name:        "STAR Registrations Index",
		Method:      "GET",
		Pattern:     "/star/registrations",
		HandlerFunc: RegistrationsList,
	},
	Route{
		Name:        "RegistrationProgress",
		Method:      "GET",
		Pattern:     "/star/registration/{registrationId}",
		HandlerFunc: RegistrationProgress,
	},
	Route{
		Name:        "STAR Registration Request",
		Method:      "POST",
		Pattern:     NewRegistrationPath,
		HandlerFunc: CreateNewRegistration,
	},
}
