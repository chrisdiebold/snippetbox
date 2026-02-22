package main

import (
	"log"
	"net/http"
)

func main() {
	// a serveMux is a router. Stores a mapping between URLs routing patterns
	// for the application and the corresponding handlers (controllers)
	mux := http.NewServeMux()
	// serve static files such as css, js, and images
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	// using {$} to prevent subtree path patterns from acting like they have a
	// wildcard at the end - ends the pattern
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Println("starting server on :4000")
	// starts a web server. If this returns an err we use the log.Fatal() function to log the
	// error message and terminate the program.
	// Note: any error returned by http.ListenAndServe() is always non-nil
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
