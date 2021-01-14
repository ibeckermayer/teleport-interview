package handlers

import (
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
)

// LogoutHandler handles calls to "/api/logout". Implements HandlerWithSession
type LogoutHandler struct {
	sm *auth.SessionManager
}

// NewLogoutHandler creates a new LogoutHandler
func NewLogoutHandler(sm *auth.SessionManager) *LogoutHandler {
	return &LogoutHandler{sm}
}

func (lh *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Couldn't find session. Log that this happened and return a 204
	// to alert the client to behave as if the user is logged out.
	if !lh.sm.DeleteSession(r.Context()) {
		log.Println("logout attempted but could not find session")
	}

	w.WriteHeader(http.StatusNoContent)
}
