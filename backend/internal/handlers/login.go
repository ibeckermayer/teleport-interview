package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
)

// LoginHandler handles calls to "/api/login". Implements http.Handler
type LoginHandler struct {
	sm *auth.SessionManager
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(sm *auth.SessionManager) *LoginHandler {
	return &LoginHandler{sm}
}

type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	Msg       string         `json:"msg"`
	SessionID auth.SessionID `json:"sessionID"`
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
	if body.Email == "admin@goteleport.com" && body.Password == "admin@goteleport.com" {
		// TODO: "0" should become a real uuid
		sessionID, err := lh.sm.CreateSession("0")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		lrb := loginResponseBody{"Sign in succeeded", sessionID}
		json.NewEncoder(w).Encode(lrb)
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}

	return
}
