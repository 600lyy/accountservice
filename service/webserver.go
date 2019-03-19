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
	http.Handle("/", r) 
	// http.HandleFunc("/accounts", routes[0].HandlerFunc)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		//log.Fatal("An error occurred starting HTTP listenser at port " + port)
		log.Fatal("Error: " + err.Error())
	}
}