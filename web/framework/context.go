package framework

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const userKey = "$_USER_KEY_$"

// Context type
type Context interface {
	Session() *Session
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Params() httprouter.Params
	User() *User
	Redirect(path string, code int)
	SetUser(*User)
}

// User struct
type User struct {
	ID   string
	Name string
}

type cxt struct {
	r *http.Request
	w http.ResponseWriter
	s *Session
	p httprouter.Params
}

func (c *cxt) Request() *http.Request {
	return c.r
}

func (c *cxt) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *cxt) User() *User {
	u, ok := c.s.Get(userKey).(*User)
	if !ok {
		return nil
	}

	return u
}

func (c *cxt) Session() *Session {
	return c.s
}

func (c *cxt) Redirect(url string, code int) {
	http.Redirect(c.w, c.r, url, code)
}

func (c *cxt) SetUser(user *User) {
	c.s.Set(userKey, user)
}

func (c *cxt) Params() httprouter.Params {
	return c.p
}
