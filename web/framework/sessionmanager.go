package framework

import (
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SessionManager stores session info
type SessionManager interface {
	FetchOrCreate(Context) *Session
	Abandon(Context)
}

// CookieExpireDelete may be set on Cookie.Expire for expiring the given cookie.
var CookieExpireDelete = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

type inMemStore struct {
	sync.RWMutex
	data map[SessionID]*Session
}

// NewInMemorySessionStore creates a new InMemory SessionStore
func NewInMemorySessionStore() SessionManager {
	return &inMemStore{
		data: make(map[SessionID]*Session),
	}
}

func (m *inMemStore) FetchOrCreate(cxt Context) *Session {
	r := cxt.Request()
	w := cxt.ResponseWriter()
	c, err := r.Cookie(".session_id")
	var sessionID SessionID
	if err == http.ErrNoCookie {
		sessionID = SessionID(uuid.New().String())
	}
	if c != nil {
		sessionID = SessionID(c.Value)
	}
	m.Lock()
	defer m.Unlock()
	session := m.data[sessionID]
	if session == nil {
		session = &Session{id: sessionID, data: make(Map)}
		m.data[sessionID] = session
	}
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Name:     ".session_id",
		Path:     "/",
		Value:    string(sessionID),
	})
	return session
}

func (m *inMemStore) Abandon(context Context) {
	m.Lock()
	defer m.Unlock()

	http.SetCookie(context.ResponseWriter(), &http.Cookie{
		HttpOnly: true,
		Secure:   context.Request().TLS != nil,
		Name:     ".session_id",
		Path:     "/",
		Value:    "",
		Expires:  CookieExpireDelete,
		MaxAge:   -1,
	})
	delete(m.data, context.Session().ID())
}
