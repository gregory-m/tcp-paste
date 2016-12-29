package server

import "net"

// TCP is simple TCP server it invokes handler on each tcp connection
type TCP struct {
	Addr    *net.TCPAddr
	Handler TCPHandler

	listener *net.TCPListener
	quit     chan bool
}

// Start server and listen for connection. Note its blocking call
func (s *TCP) Start() error {
	ln, err := net.ListenTCP("tcp", s.Addr)
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
		go s.Handler.Handle(conn.(*net.TCPConn)) //TODO report error
	}
}

// Stop server
// TODO stop it gradually
func (s *TCP) Stop() error {
	s.listener.Close()
	return nil
}
