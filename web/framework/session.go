package framework

import (
	"sync"
)

// SessionID unique session identifier
type SessionID string

// Session represents web session
type Session struct {
	sync.RWMutex
	id   SessionID
	data map[string]interface{}
}

// Set the key value pair specified to the session
func (s *Session) Set(key string, val interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = val
}

// Get the value for the specified key
func (s *Session) Get(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.data[key]
}

// ID returns the SessionID
func (s *Session) ID() SessionID {
	return s.id
}
