package shared

import "net/http"

// FileProvider fetches files by name
type FileProvider interface {
	List() []string
	String(name string) string
	Open(name string) (http.File, error)
	Watch()
	Close()
}
