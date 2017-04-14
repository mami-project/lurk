package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"LURK Registrations Index",
		"GET",
		"/lurk/registrations",
		RegistrationsList,
	},
	Route{
		"RegistrationProgress",
		"GET",
		"/lurk/registration/{registrationId}",
		RegistrationProgress,
	},
	Route{
		"LURK Registration Request",
		"POST",
		"/lurk/registration",
		RegistrationCreate,
	},
}
