package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/database"
)

// LoginHandler handles calls to "/api/login". Implements http.Handler
type LoginHandler struct {
	sm *auth.SessionManager
	db *database.Database
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(sm *auth.SessionManager, db *database.Database) *LoginHandler {
	return &LoginHandler{sm, db}
}

type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	SessionID auth.SessionID `json:"sessionID"`
	Plan      string         `json:"plan"`
}

// Handles user login
func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody

	err := decodeJSONBody(w, r, &body)
	if err != nil {
		handleJSONdecodeError(w, err)
		return
	}

	account, err := lh.db.GetAccount(body.Email)

	// Handle errors from attempting to retrieve the account from the database
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			// No record with the given email address exists
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Account retrieved, check password
	if !auth.CheckPasswordHash(body.Password, account.PasswordHash) {
		// Invalid password, unauthorized
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Valid password, create new session
	session, err := lh.sm.CreateSession(account)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(loginResponseBody{session.SessionID, session.Account.Plan})
}
