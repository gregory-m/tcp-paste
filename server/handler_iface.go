package server

import "net"

// TCPHandler used by TCP server to handle connections
type TCPHandler interface {
	// Handle should close conn
	Handle(conn *net.TCPConn) error
}
