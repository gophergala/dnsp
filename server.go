// Package dnsp contains a simple DNS proxy.
package dnsp

import "net"

// Server implements a DNS server.
type Server struct {
	conn *net.UDPConn
}

// NewServer creates a new Server with the given options.
func NewServer(o Options) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", o.Bind)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn: conn,
	}, nil
}

// Start runs the server
func (s *Server) Start() error {
	select {} // TODO
}

// Stop stops the server, closing any kernel buffers.
func (s *Server) Stop() error {
	return s.conn.Close()
}
