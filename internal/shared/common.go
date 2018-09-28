package common

// FileProvider fetches files by name
type FileProvider interface {
	List() []string
	String(name string) string
}
