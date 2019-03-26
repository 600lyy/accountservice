package service

import (
	"net/http"
	"log"
)

// StartWebServer starts a HTTP server at a
// specified port
func StartWebServer(port string) {
	log.Println("Starting HTTP server at " + port)
	r := NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler) // display 404 html
	http.Handle("/", r) 
	// http.HandleFunc("/accounts", routes[0].HandlerFunc)
	err := http.ListenAndServe(":" + port, r)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
}