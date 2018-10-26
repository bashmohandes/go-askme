package framework

import (
	"container/heap"
	"log"
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

type sessionheap []*heapitem

type heapitem struct {
	sid   SessionID
	index int
}

// CookieExpireDelete may be set on Cookie.Expire for expiring the given cookie.
var CookieExpireDelete = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

type inMemStore struct {
	sync.RWMutex
	data map[SessionID]*Session
	pq   sessionheap
}

// NewInMemorySessionStore creates a new InMemory SessionStore
func NewInMemorySessionStore() SessionManager {
	mem := &inMemStore{
		data: make(map[SessionID]*Session),
	}

	heap.Init(mem)
	return mem
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
	expires := 2 * time.Minute
	if session == nil {
		session = &Session{id: sessionID, data: make(Map), expires: time.Now().Add(expires)}
		m.data[sessionID] = session
		heap.Push(m, &heapitem{sid: sessionID})
		if len(m.pq) == 1 {
			time.AfterFunc(expires, m.GC)
		}
	}
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Name:     ".session_id",
		Path:     "/",
		Value:    string(sessionID),
		Expires:  session.expires,
		MaxAge:   int(expires.Seconds()),
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
	for i := range m.pq {
		if m.pq[i].sid == context.Session().ID() {
			heap.Remove(m, i)
			return
		}
	}
}

func (m *inMemStore) Len() int {
	return len(m.pq)
}

func (m *inMemStore) Less(i, j int) bool {
	si := m.data[m.pq[i].sid]
	sj := m.data[m.pq[j].sid]

	return si.expires.Before(sj.expires)
}

func (m *inMemStore) Swap(i, j int) {
	m.pq[i], m.pq[j] = m.pq[j], m.pq[i]
	m.pq[i].index = i
	m.pq[j].index = j
}

func (m *inMemStore) Push(x interface{}) {
	n := len(m.pq)
	item := x.(*heapitem)
	item.index = n
	m.pq = append(m.pq, item)
}

func (m *inMemStore) Pop() interface{} {
	old := m.pq
	n := len(old)
	x := old[n-1]
	x.index = -1
	m.pq = old[0 : n-1]
	return x
}

func (m *inMemStore) GC() {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	log.Printf("GC called len(m.pq)=%d, len(m.data)=%d", len(m.pq), len(m.data))
	if len(m.pq) != 0 {
		topMostTime := m.data[m.pq[0].sid].expires
		if topMostTime.Before(time.Now()) { // expired
			log.Printf("GC cleaning session=%s", m.pq[0].sid)
			delete(m.data, m.pq[0].sid)
			heap.Pop(m)
		}
	}
	nextCheck := 10 * time.Second
	if len(m.pq) > 0 {
		nextCheckTop := m.data[m.pq[0].sid].expires.Sub(time.Now())
		if nextCheckTop > nextCheck {
			nextCheck = nextCheckTop
		}
	}
	log.Printf("GC scheduled after %v seconds", nextCheck.Seconds())
	time.AfterFunc(nextCheck, m.GC)
}
