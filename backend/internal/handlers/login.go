package handlers

import (
	"encoding/json"
	"net/http"
)

// LoginHandler handles calls to "/api/login". Implements http.Handler
type LoginHandler struct {
	// TODO: Pass db/session manager pointers through from server.Server
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	Msg string `json:"msg"`
}

// Just for testing purposes, TODO: should be deleted
func fakeAuthLogic(body *loginRequestBody, w http.ResponseWriter, r *http.Request) {
	if body.Email == "admin@goteleport.com" && body.Password == "admin@goteleport.com" {
		lrb := loginResponseBody{"Sign in succeeded"}
		json.NewEncoder(w).Encode(lrb)
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

// Handles user login
func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody

	err := decodeJSONBody(w, r, &body)
	if err != nil {
		handleJSONdecodeError(w, err)
		return
	}

	// TODO: remove
	fakeAuthLogic(&body, w, r)
	return
}
