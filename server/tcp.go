package server

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path"
)

// TCP represents TCP server
type TCP struct {
	HostName   string
	StorageDir string
	Host       *net.TCPAddr

	listener *net.TCPListener
	quit     chan bool
}

// Start server
func (s *TCP) Start() error {
	ln, err := net.ListenTCP("tcp", s.Host)
	if err != nil {
		return err
	}
	s.listener = ln

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return nil
			default:
				//Without default select will block
			}

			// TODO report error
			continue
		}
		go s.handleConn(conn.(*net.TCPConn)) //TODO report error
	}
}

func (s *TCP) handleConn(conn *net.TCPConn) error {
	fmt.Printf("%s connected\n", conn.RemoteAddr().String())
	defer conn.Close()
	fileName := randomString(10)

	file, err := os.Create(path.Join(s.StorageDir, fileName))
	defer file.Close()
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("Sorry can't create file: %s", err)))
		return err
	}

	writed, err := io.Copy(file, conn)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("Sorry write to file: %s", err)))
		return err
	}

	u := &url.URL{
		Scheme: "http",
		Host:   s.HostName,
		Path:   fileServerPrefix,
	}
	u.Path = path.Join(u.Path, fileName)

	fmt.Printf("Saved %d bytes for %s at %s\n", writed, conn.RemoteAddr().String(),
		u.String())

	_, err = conn.Write([]byte(fmt.Sprintf("Results saved at: %s\n", u.String())))
	if err != nil {
		fmt.Printf("can't write to connection: %s\n", err)
	}

	conn.Close()
	return nil
}

//Stop server
func (s *TCP) Stop() error {
	s.listener.Close()
	return nil
}
