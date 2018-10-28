package framework

import (
	"container/heap"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestInMemorySessionHeap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	mem := &inMemStore{
		data:            make(map[SessionID]*Session),
		sessionCookie:   ".session_id",
		sessionLifetime: time.Second,
	}
	heap.Init(mem)

	for _, i := range []*Session{newSession(), newSession(), newSession(), newSession()} {
		fmt.Printf("S: %s, expires: %s\n", i.id, i.expires)
		heap.Push(mem, i)
	}
	fmt.Printf("Session with min %s, time %s\n", mem.min().id, mem.min().expires)
	fmt.Println("Printing in order")
	for len(mem.pq) > 0 {
		s := heap.Pop(mem).(*Session)
		fmt.Printf("S: %s, len(mem): %d\n", s.id, len(mem.data))
	}
}

func newSession() *Session {
	r, _ := time.ParseDuration(fmt.Sprintf("%ds", rand.Intn(1000)))
	fmt.Println(r)
	return &Session{
		id:      SessionID(uuid.New().String()),
		expires: time.Now().Add(r),
		data:    make(map[string]interface{}),
	}
}
