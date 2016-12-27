package server

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
)

func TestOutput(t *testing.T) {
	testHost := "example.test"

	ts := newTestSerevr(testHost)
	defer ts.cleenUP()

	outputURL := writeAndGetURL(ts, []byte("testa"), t)

	if outputURL.Host != "example.test" {
		t.Errorf("Wrong host in returned URl\nExpected: %q\nGot:      %q\n", testHost, outputURL.Host)
	}

	if !strings.HasPrefix(outputURL.Path, fileServerPrefix) {
		t.Errorf("Wrong prefix in returned URl\nExpected: %q\nGot:      %q\n", "/"+fileServerPrefix, outputURL.Path)
	}
}

func TestFileContent(t *testing.T) {
	testInput := []byte("testa 123")

	ts := newTestSerevr("example.com")
	defer ts.cleenUP()

	outputURL := writeAndGetURL(ts, testInput, t)

	fileName := path.Base(outputURL.Path)

	output, err := ioutil.ReadFile(path.Join(ts.StorageDir, fileName))
	if err != nil {
		t.Fatalf("Can't read file: %s", err)
	}

	if !reflect.DeepEqual(output, testInput) {
		t.Errorf("File content dosn't match\nExpected: %q\nGot:      %q\n", testInput, output)
	}
}

func newTestSerevr(hostname string) *TCP {
	service := ":4343"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		panic(err)
	}

	tmpDir, err := ioutil.TempDir("", "fiche-go-tests")
	if err != nil {
		panic(err)
	}

	return &TCP{HostName: hostname, Host: tcpAddr, StorageDir: tmpDir}
}

func newTestConn(s *TCP, t *testing.T) *net.TCPConn {

	go s.Start()
	time.Sleep(1 * time.Millisecond)

	c, err := net.Dial("tcp", "127.0.0.1:4343")
	if err != nil {
		t.Fatalf("Can't dial to server: %s", err)
	}

	return c.(*net.TCPConn)
}

func writeAndGetURL(s *TCP, input []byte, t *testing.T) *url.URL {
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

func parseOutput(in string) (*url.URL, error) {
	parts := strings.Split(in, " ")
	outputURL, err := url.Parse(strings.TrimSuffix(parts[len(parts)-1], "\n"))

	if err != nil {
		panic(err)
	}

	return outputURL, nil
}

func (s *TCP) cleenUP() {
	s.Stop()
	os.RemoveAll(s.StorageDir)
}
