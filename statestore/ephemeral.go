package statestore

import (
	"sync"
	"time"
)

const (
	maxEphemeralSessions   = 10_000
	ephemeralTTL           = 2 * time.Hour
	ephemeralEvictInterval = 15 * time.Minute
)

// EphemeralState holds transient UI state for a session.
// It is never used as input to durable state transitions.
type EphemeralState struct {
	CurrentPage string
	FlashMsg    string
}

type ephemeralEntry struct {
	state        *EphemeralState
	lastAccessed time.Time
}

type ephemeralStore struct {
	mu       sync.Mutex
	sessions map[string]*ephemeralEntry
}

func newEphemeralStore() *ephemeralStore {
	s := &ephemeralStore{sessions: make(map[string]*ephemeralEntry)}
	go s.runEviction()
	return s
}

func (e *ephemeralStore) runEviction() {
	t := time.NewTicker(ephemeralEvictInterval)
	for range t.C {
		cutoff := time.Now().Add(-ephemeralTTL)
		e.mu.Lock()
		for k, v := range e.sessions {
			if v.lastAccessed.Before(cutoff) {
				delete(e.sessions, k)
			}
		}
		e.mu.Unlock()
	}
}

func (e *ephemeralStore) Get(sessionID string) *EphemeralState {
	e.mu.Lock()
	defer e.mu.Unlock()
	entry, ok := e.sessions[sessionID]
	if !ok {
		return &EphemeralState{}
	}
	entry.lastAccessed = time.Now()
	return entry.state
}

func (e *ephemeralStore) Set(sessionID string, state *EphemeralState) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.sessions[sessionID]; !exists && len(e.sessions) >= maxEphemeralSessions {
		e.evictOldest()
	}
	e.sessions[sessionID] = &ephemeralEntry{state: state, lastAccessed: time.Now()}
}

// evictOldest removes the least-recently-accessed session. Must be called with e.mu held.
func (e *ephemeralStore) evictOldest() {
	var oldest string
	var oldestTime time.Time
	for k, v := range e.sessions {
		if oldest == "" || v.lastAccessed.Before(oldestTime) {
			oldest = k
			oldestTime = v.lastAccessed
		}
	}
	if oldest != "" {
		delete(e.sessions, oldest)
	}
}

func (e *ephemeralStore) Delete(sessionID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.sessions, sessionID)
}
