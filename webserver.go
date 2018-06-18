package main

/*
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host
Navigating to http://localhost:8100 will display the index.html or directory
listing file.
*/

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Serve serves a very basic static page
func serve() error {

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
	return nil
}
