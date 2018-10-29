package middleware

import (
	"log"

	"github.com/bashmohandes/go-askme/web/framework"
)

// RequestLogger logs all requests
type requestLogger struct {
}

// Run the request logger middleware
func (l *requestLogger) Run(cxt framework.Context) bool {
	r := cxt.Request()
	log.Printf("Request Recvd, Method: %s, path: %s, req size: %d", r.Method, r.RequestURI, r.ContentLength)
	return true
}

// NewRequestLogger creates a new logger
func NewRequestLogger() framework.Middleware {
	return &requestLogger{}
}
