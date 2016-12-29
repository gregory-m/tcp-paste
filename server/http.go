package server

import (
	"net"
	"net/http"
	"strings"
)

// HTTP represents TCP server
type HTTP struct {
	StorageDir string
	Host       *net.TCPAddr
	Prefix     string
}

// Start http server
func (s *HTTP) Start() error {
	http.Handle(FileServerPrefix, serveStatic(s.StorageDir))
	http.HandleFunc("/", homePage)

	return http.ListenAndServe(s.Host.String(), nil)
}

func homePage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is tcp-paste server.\n"))
}

func serveStatic(dir string) http.Handler {
	return noDirListing(http.StripPrefix(FileServerPrefix, http.FileServer(http.Dir(dir))))
}

func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
