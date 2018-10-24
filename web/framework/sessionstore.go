package framework

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

// SessionStore stores session info
type SessionStore interface {
	Get(SessionID) (*Session, error)
	FetchOrCreate(http.ResponseWriter, *http.Request) *Session
	Abandon(SessionID)
}

type inMemStore struct {
	sync.RWMutex
	data map[SessionID]*Session
}

// NewInMemorySessionStore creates a new InMemory SessionStore
func NewInMemorySessionStore() SessionStore {
	return &inMemStore{
		data: make(map[SessionID]*Session),
	}
}

func (m *inMemStore) FetchOrCreate(w http.ResponseWriter, r *http.Request) *Session {
	m.Lock()
	defer m.Unlock()
	c, err := r.Cookie(".auth")
	var sessionID SessionID
	if err == http.ErrNoCookie {
		sessionID = SessionID(uuid.New().String())
	}

	if c != nil {
		sessionID = SessionID(c.Value)
	}

	session := &Session{id: sessionID, data: make(map[string]interface{})}
	m.data[sessionID] = session
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Name:     ".auth",
		Path:     "/",
		Value:    string(sessionID),
	})
	return session
}

func (m *inMemStore) Get(id SessionID) (*Session, error) {
	m.RLock()
	defer m.RUnlock()
	s, ok := m.data[id]
	if !ok {
		return nil, fmt.Errorf("Session not found")
	}

	return s, nil
}

func (m *inMemStore) Abandon(id SessionID) {
	m.Lock()
	defer m.Unlock()

	delete(m.data, id)
}
