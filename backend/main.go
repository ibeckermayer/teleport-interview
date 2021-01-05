package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ibeckermayer/teleport-interview/backend/internal/handlers"
)

func main() {
	router := mux.NewRouter()

	router.Handle("/api/login", http.HandlerFunc(handlers.Login)).Methods("POST")

	// NOTE: It's important that this handler be registered after the other handlers, or else
	// all routes return a 404 (at least in development). TODO: figure out why this is the case.
	router.PathPrefix("/").Handler(handlers.NewSpaHandler("../frontend", "index.html"))

	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServeTLS(":8000", "../certs/localhost.crt", "../certs/localhost.key", router))
}
