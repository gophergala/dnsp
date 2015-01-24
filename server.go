// Package dnsp contains a simple DNS proxy.
package dnsp

// Server implements a DNS server.
type Server struct{}

// NewServer creates a new Server with the given options.
func NewServer(o Options) (*Server, error) {
	return &Server{}, nil
}

// Start runs the server
func (s *Server) Start() error {
	select {} // TODO
	return nil
}

// Stop stops the server, closing any kernel buffers.
func (s *Server) Stop() error {
	return nil
}
