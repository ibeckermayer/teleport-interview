package auth

import (
	"sync"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

type session struct {
	account *model.Account
	expires time.Time
}

// SessionManager is an in-memory session store
type SessionManager struct {
	store   map[SessionID]*session
	timeout time.Duration // absolute timeout for individual sessions
	mtx     sync.RWMutex
}

// NewSessionManager creates a new *SessionManager
func NewSessionManager(timeout time.Duration) *SessionManager {
	return &SessionManager{store: make(map[SessionID]*session), timeout: timeout}
}

// CreateSession creates a new session in the SessionManager's store, indexed by a
// new randomly generated SessionID, and expiring sm.timeout from the time it's created.
// It will return an error if the system's secure random number generator fails to function
// correctly, in which case the caller should not continue.
func (sm *SessionManager) CreateSession(account *model.Account) (SessionID, error) {
	sid, err := newSessionID()
	if err != nil {
		return "", err
	}

	// TODO: Do we really need an RLock here to read sm.timeout?
	sm.mtx.RLock()
	s := &session{account, time.Now().Add(sm.timeout)}
	sm.mtx.RUnlock()

	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.store[sid] = s

	return sid, nil
}
