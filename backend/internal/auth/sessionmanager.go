package auth

import (
	"sync"
	"time"
)

// AccountID is a fakeiot account's uuid
type AccountID string

type session struct {
	accountID AccountID
	expires   time.Time
}

// SessionManager is an in-memory session store
type SessionManager struct {
	store   map[SessionID]session
	timeout time.Duration // Read only, does not need synchronization
	mtx     sync.RWMutex
}

// NewSessionManager creates a new *SessionManager
func NewSessionManager(timeout time.Duration) *SessionManager {
	return &SessionManager{store: make(map[SessionID]session), timeout: timeout}
}

// CreateSession creates a new session in the SessionManager's store, indexed by a
// new randomly generated SessionID, and expiring sm.timeout from the time it's created.
// It will return an error if the system's secure random number generator fails to function
// correctly, in which case the caller should not continue.
func (sm *SessionManager) CreateSession(accountID AccountID) (SessionID, error) {
	sid, err := newSessionID()
	if err != nil {
		return "", err
	}

	s := session{accountID, time.Now().Add(sm.timeout)}

	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.store[sid] = s

	return sid, nil
}
