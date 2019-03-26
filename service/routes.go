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
		Pattern:     "/accounts/{username}",
		HandlerFunc: GetAccount,
	},

	Route{
		"HealthCheck",
		"GET",
		"/health",
		HealthCheck,
	},

	Route{
		Name:        "Regisration",
		Method:      "POST",
		Pattern:     "/register",
		HandlerFunc: CreateAccount,
	},

	Route{
		Name:        "Regisration",
		Method:      "GET",
		Pattern:     "/register/index",
		HandlerFunc: CreateAccount,
	},

	Route{
		Name:        "GetAllDemoAccounts",
		Method:      "GET",
		Pattern:     "/accounts	",
		HandlerFunc: GetAllDemoAccounts,
	},

	Route{
		Name:        "Login",
		Method:      "POST",
		Pattern:     "/login",
		HandlerFunc: UserLogin,
	},

	Route{
		Name:        "Login",
		Method:      "GET",
		Pattern:     "/login/index",
		HandlerFunc: UserLogin,
	},

}
