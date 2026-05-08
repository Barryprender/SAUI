package statestore

import "sync"

// EphemeralState holds transient UI state for a session.
// It is never used as input to durable state transitions.
type EphemeralState struct {
	CurrentPage string
	FlashMsg    string
}

type ephemeralStore struct {
	mu      sync.RWMutex
	sessions map[string]*EphemeralState
}

func newEphemeralStore() *ephemeralStore {
	return &ephemeralStore{sessions: make(map[string]*EphemeralState)}
}

func (e *ephemeralStore) Get(sessionID string) *EphemeralState {
	e.mu.RLock()
	defer e.mu.RUnlock()
	s, ok := e.sessions[sessionID]
	if !ok {
		return &EphemeralState{}
	}
	return s
}

func (e *ephemeralStore) Set(sessionID string, state *EphemeralState) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.sessions[sessionID] = state
}

func (e *ephemeralStore) Delete(sessionID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.sessions, sessionID)
}
