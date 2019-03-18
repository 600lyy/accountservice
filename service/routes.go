package service

import (
	"net/http"
)

// Route defines a single route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a slice including many routes
type Routes []Route

var routes = Routes{

	Route{
		Name:        "GetAccount",
		Method:      "GET",
		Pattern:     "/accounts/{accountId}",
		HandlerFunc: GetAccount,
	},

	Route{
		"HealthCheck",
		"GET",
		"/health",
		HealthCheck,
	},

	Route{
		Name:        "Login",
		Method:      "GET",
		Pattern:     "/login",
		HandlerFunc: Login,
	},

	Route{
		Name:        "SignUp",
		Method:      "POST",
		Pattern:     "/signup",
		HandlerFunc: SignUp,
	},
}
