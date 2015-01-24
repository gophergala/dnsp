// Package dnsp contains a simple DNS proxy.
package dnsp

import (
	"log"
	"net"
)

// Server implements a DNS server.
type Server struct {
	conn *net.UDPConn

	blacklist map[string]bool
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
		conn:      conn,
		blacklist: map[string]bool{},
	}, nil
}

// Start runs the server
func (s *Server) Start() error {
	for {
		b, oob := make([]byte, 1024), make([]byte, 64)
		n, oobn, flags, addr, err := s.conn.ReadMsgUDP(b, oob)
		if err != nil {
			log.Printf("error=conn_read_msg details=%q", err)
			continue
		}
		log.Printf("debug=read flags=%d addr=%s b=%q oob=%q", flags, addr, b[:n], oob[:oobn])
	}
}

// Stop stops the server, closing any kernel buffers.
func (s *Server) Stop() error {
	return s.conn.Close()
}

func (s *Server) Addr() net.Addr {
	return s.conn.LocalAddr()
}
