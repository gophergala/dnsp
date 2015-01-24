package dnsp

import "github.com/miekg/dns"

// Block adds a host to the blocklist.
func (s *Server) Block(host string) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}

	if !s.IsBlocked(host) { // avoid duplicates
		s.blacklist[host] = true
	}
}

// IsBlocked returns whether a hostname is blocked.
func (s *Server) IsBlocked(host string) bool {
	return s.blacklist[host]
}

func (s *Server) filter(qs []dns.Question) []dns.Question {
	result := []dns.Question{}
	for _, q := range qs {
		if !s.IsBlocked(q.Name) {
			result = append(result, q)
		}
	}
	return result
}
