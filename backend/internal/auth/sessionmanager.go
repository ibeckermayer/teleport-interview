package auth

import (
	"sync"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

// Session is an individual user's session
type Session struct {
	SessionID SessionID
	Account   model.Account
	Expires   time.Time
}

// SessionManager is an in-memory session store. The app should only ever create one
// of these and pass it around as a pointer
type SessionManager struct {
	store   map[SessionID]Session
	timeout time.Duration // absolute timeout for individual sessions
	mtx     sync.RWMutex
}

// NewSessionManager creates a new *SessionManager
func NewSessionManager(timeout time.Duration) *SessionManager {
	return &SessionManager{store: make(map[SessionID]Session), timeout: timeout}
}

// CreateSession creates a new session in the SessionManager's store, indexed by a
// new randomly generated SessionID, and expiring sm.timeout from the time it's created.
// It will return an error if the system's secure random number generator fails to function correctly.
func (sm *SessionManager) CreateSession(account model.Account) (Session, error) {
	sid, err := newSessionID()
	if err != nil {
		return Session{}, err
	}

	s := Session{sid, account, time.Now().Add(sm.timeout)}

	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.store[sid] = s

	return s, nil
}

// DeleteSession deletes a session from the session manager. Returns true if the session
// was found and deleted, or false if the session wasn't found
func (sm *SessionManager) DeleteSession(sid SessionID) bool {
	// Check whether the session exists and return false if it doesn't
	sm.mtx.RLock()
	_, ok := sm.store[sid]
	sm.mtx.RUnlock()
	if !ok {
		return ok
	}
	// Session does exist, delete it
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	delete(sm.store, sid)

	return ok
}
