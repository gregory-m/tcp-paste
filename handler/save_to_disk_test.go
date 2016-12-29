package handler

import (
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gregory-m/tcp-paste/server"
)

func TestOutput(t *testing.T) {
	testHost := "example.test"
	ts := newTestSerevr(testHost)
	defer cleenUP(ts)

	outputURL := writeAndGetURL(ts, []byte("testa"), t)

	if outputURL.Host != "example.test" {
		t.Errorf("Wrong host in returned URl\nExpected: %q\nGot:      %q\n", testHost, outputURL.Host)
	}

	if !strings.HasPrefix(outputURL.Path, server.FileServerPrefix) {
		t.Errorf("Wrong prefix in returned URl\nExpected: %q\nGot:      %q\n", server.FileServerPrefix, outputURL.Path)
	}
}

func TestFileContent(t *testing.T) {
	testInput := []byte("testa 123")

	ts := newTestSerevr("example.com")
	defer cleenUP(ts)

	outputURL := writeAndGetURL(ts, testInput, t)

	fileName := path.Base(outputURL.Path)

	output, err := ioutil.ReadFile(path.Join(ts.Handler.(*SaveToDisk).StorageDir, fileName))
	if err != nil {
		t.Fatalf("Can't read file: %s", err)
	}

	if !reflect.DeepEqual(output, testInput) {
		t.Errorf("File content dosn't match\nExpected: %q\nGot:      %q\n", testInput, output)
	}
}

func newTestHandler(host string) server.TCPHandler {
	tmpDir, err := ioutil.TempDir("", "fiche-go-tests")
	if err != nil {
		panic(err)
	}

	return &SaveToDisk{
		HostName:   host,
		StorageDir: tmpDir,
		Prefix:     server.FileServerPrefix,
	}
}

func newTestSerevr(hostName string) *server.TCP {
	h := newTestHandler(hostName)
	a, err := net.ResolveTCPAddr("tcp", ":4343")
	if err != nil {
		panic(err)
	}

	s := &server.TCP{
		Handler: h,
		Addr:    a,
	}

	return s
}

func cleenUP(s *server.TCP) {
	s.Stop()
	os.RemoveAll(s.Handler.(*SaveToDisk).StorageDir)
}

func writeAndGetURL(s *server.TCP, input []byte, t *testing.T) *url.URL {
	conn := newTestConn(s, t)

	_, err := conn.Write(input)
	if err != nil {
		t.Fatalf("Can't write: %s", err)
	}
	conn.CloseWrite()
	if err != nil {
		t.Fatalf("Can't CloseWrite connection: %s", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Can't read from connection: %s", err)
	}
	out := string(buf[:n])

	url, err := parseOutput(out)
	if err != nil {
		t.Fatalf("Can't parse URL: %s", err)
	}

	return url
}

func newTestConn(s *server.TCP, t *testing.T) *net.TCPConn {

	go s.Start()
	time.Sleep(1 * time.Millisecond)

	c, err := net.Dial("tcp", "127.0.0.1:4343")
	if err != nil {
		t.Fatalf("Can't dial to server: %s", err)
	}

	return c.(*net.TCPConn)
}

func parseOutput(in string) (*url.URL, error) {
	parts := strings.Split(in, " ")
	outputURL, err := url.Parse(strings.TrimSuffix(parts[len(parts)-1], "\n"))

	if err != nil {
		panic(err)
	}

	return outputURL, nil
}
