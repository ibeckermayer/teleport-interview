package handlers

import (
	"errors"
	"log"
	"net/http"
)

type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var v http.HandlerFunc

// Login handles user logging in
func Login(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody

	err := decodeJSONBody(w, r, &body)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			log.Println(mr.Error())
			http.Error(w, http.StatusText(mr.status), mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// TODO: pwd should be saved in db and bcrypt-ified
	if body.Email == "admin@goteleport.com" && body.Password == "admin@goteleport.com" {
		w.Write([]byte("Sign in succeeded"))
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
