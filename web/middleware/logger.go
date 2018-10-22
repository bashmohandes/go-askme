package middleware

import (
	"log"
	"net/http"

	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/julienschmidt/httprouter"
)

// RequestLogger logs all requests
type requestLogger struct {
}

// Run the request logger middleware
func (l *requestLogger) Run(w http.ResponseWriter, r *http.Request, params httprouter.Params) bool {
	log.Printf("Request Recvd, Method: %s, path: %s, req size: %d", r.Method, r.RequestURI, r.ContentLength)
	return true
}

// NewRequestLogger creates a new logger
func NewRequestLogger() framework.Middleware {
	return &requestLogger{}
}
