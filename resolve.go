package dnsp

import "github.com/miekg/dns"

// IsBlocked returns true if a DNS query is blocked.
func (s *Server) IsBlocked(*dns.Msg) bool {
	// TODO: consult s.IsHostBlocked()
	return false
}

// IsHostBlocked returns true if a hostname is blocked.
func (s *Server) IsHostBlocked(host string) bool {
	// TODO: consult the local block list.
	return host == "blocked.net"
}
