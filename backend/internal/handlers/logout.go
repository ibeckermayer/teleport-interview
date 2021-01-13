package handlers

import (
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
)

// LogoutHandler handles calls to "/api/logout". Implements http.Handler
type LogoutHandler struct {
	sm *auth.SessionManager
}

// NewLogoutHandler creates a new LogoutHandler
func NewLogoutHandler(sm *auth.SessionManager) *LogoutHandler {
	return &LogoutHandler{sm}
}

func (lh *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)

	// Couldn't retrieve token. Log that the error occurred but return a 204
	// to alert the client to behave as if the user is logged out.
	if sessionID == "" {
		log.Println("logout attempted but there was an empty or improperly formatted Authorization header, no session deleted")
	}

	// Couldn't find session. Log that this happened and return a 204
	// to alert the client to behave as if the user is logged out.
	if !lh.sm.DeleteSession(sessionID) {
		log.Printf("logout attempted but could not find session %v", sessionID)
	}

	w.WriteHeader(http.StatusNoContent)
}
