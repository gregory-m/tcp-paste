package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestStaticServe(t *testing.T) {
	tf := createTestFile()
	defer os.Remove(tf)

	req, _ := http.NewRequest("GET", path.Join(FileServerPrefix, path.Base(tf)), nil)
	w := httptest.NewRecorder()

	serveStatic("/tmp").ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Return %v on file request", w.Code)
	}
}

func TestNotDirListings(t *testing.T) {
	tf := createTestFile()
	defer os.Remove(tf)

	req, _ := http.NewRequest("GET", FileServerPrefix+"/", nil)
	w := httptest.NewRecorder()

	serveStatic("/tmp").ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Return %v on directtory listing", w.Code)
	}
}

func createTestFile() string {
	f, err := ioutil.TempFile("/tmp", "bla-bla")
	if err != nil {
		panic(err)
	}
	f.Close()
	return f.Name()
}
