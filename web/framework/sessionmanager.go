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

type sessionheap []SessionID

// CookieExpireDelete may be set on Cookie.Expire for expiring the given cookie.
var CookieExpireDelete = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

type sessionManager struct {
	sync.Mutex
	data            map[SessionID]*Session
	pq              sessionheap
	sessionLifetime time.Duration
	sessionCookie   string
	timer           *time.Timer
}

// NewInMemorySessionStore creates a new InMemory SessionStore
func NewInMemorySessionStore(config *Config) SessionManager {
	mem := &sessionManager{
		data:            make(map[SessionID]*Session),
		sessionCookie:   config.SessionCookie,
		sessionLifetime: config.SessionMaxLifeTime,
	}
	heap.Init(mem)
	mem.timer = time.AfterFunc(config.SessionMaxLifeTime, mem.GC)
	return mem
}

func (m *sessionManager) FetchOrCreate(cxt Context) *Session {
	r := cxt.Request()
	w := cxt.ResponseWriter()
	c, err := r.Cookie(m.sessionCookie)
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
	expires := time.Now().Add(m.sessionLifetime)
	if session == nil {
		session = &Session{id: sessionID, data: make(Map), expires: expires}
		heap.Push(m, session)
	} else {
		session.expires = expires
		heap.Fix(m, session.heapIndex)
	}
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Name:     m.sessionCookie,
		Path:     "/",
		Value:    string(sessionID),
		Expires:  session.expires,
		MaxAge:   int(m.sessionLifetime.Seconds()),
	})

	return session
}

func (m *sessionManager) Abandon(context Context) {
	m.Lock()
	defer m.Unlock()

	http.SetCookie(context.ResponseWriter(), &http.Cookie{
		HttpOnly: true,
		Secure:   context.Request().TLS != nil,
		Name:     m.sessionCookie,
		Path:     "/",
		Value:    "",
		Expires:  CookieExpireDelete,
		MaxAge:   -1,
	})
	heap.Remove(m, m.data[context.Session().ID()].heapIndex)
	delete(m.data, context.Session().ID())
}

func (m *sessionManager) Len() int {
	return len(m.pq)
}

func (m *sessionManager) Less(i, j int) bool {
	si := m.data[m.pq[i]]
	sj := m.data[m.pq[j]]

	return si.expires.Before(sj.expires)
}

func (m *sessionManager) Swap(i, j int) {
	m.pq[i], m.pq[j] = m.pq[j], m.pq[i]
	m.data[m.pq[i]].heapIndex = i
	m.data[m.pq[j]].heapIndex = j
}

func (m *sessionManager) Push(x interface{}) {
	n := len(m.pq)
	item := x.(*Session)
	m.data[item.ID()] = item
	m.data[item.ID()].heapIndex = n
	m.pq = append(m.pq, item.ID())
}

func (m *sessionManager) Pop() interface{} {
	old := m.pq
	n := len(old)
	sid := old[n-1]
	sess := m.data[sid]
	m.data[sid].heapIndex = -1
	m.pq = old[0 : n-1]
	delete(m.data, sid)
	return sess
}

func (m *sessionManager) min() *Session {
	return m.data[m.pq[0]]
}

func (m *sessionManager) GC() {
	m.Lock()
	defer m.Unlock()

	//log.Printf("GC called len(m.pq)=%d, len(m.data)=%d", len(m.pq), len(m.data))
	for len(m.pq) > 0 {
		topMostTime := m.min().expires
		//log.Printf("Session with closes expiry time is %s time %s", m.min().id, topMostTime)
		if topMostTime.Before(time.Now()) { // expired session
			delSess := heap.Pop(m).(*Session)
			log.Printf("Popped Session %s", delSess.id)
		} else {
			break
		}
	}
	nextCheck := 10 * time.Second
	if len(m.pq) > 0 {
		nextCheck = m.min().expires.Sub(time.Now())
	}
	//log.Printf("GC scheduled after %v seconds, len(timer.C)=%d", nextCheck.Seconds(), len(m.timer.C))

	m.timer.Reset(nextCheck)
}
