package server

// FileServerPrefix prefix to serve static files
var FileServerPrefix = "/f/"

// Service interface for both TCP and HTTP servers
type Service interface {
	Start() error
	Stop() error
}
