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

// GetSessionManager returns the global SessionManager
func (lh *LogoutHandler) GetSessionManager() *auth.SessionManager {
	return lh.sm
}

func (lh *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, s auth.Session) {
	// Couldn't find session. Log that this happened and return a 204
	// to alert the client to behave as if the user is logged out.
	if !lh.sm.DeleteSession(s.SessionID) {
		log.Printf("logout attempted but could not find session %v", s.SessionID)
	}

	w.WriteHeader(http.StatusNoContent)
}
