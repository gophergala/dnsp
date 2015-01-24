package dnsp

import "github.com/miekg/dns"

func (s *Server) Block(host string) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}

	if !s.IsHostBlocked(host) { // avoid duplicates
		s.blacklist[host] = true
	}
}

// IsBlocked returns whether each question in msg.Questions should be blocked.
func (s *Server) IsBlocked(msg *dns.Msg) []bool {
	results := make([]bool, len(msg.Question), len(msg.Question))
	for i, q := range msg.Question {
		results[i] = s.IsHostBlocked(q.Name)
	}
	return results
}

// IsHostBlocked returns true if a hostname is blocked.
func (s *Server) IsHostBlocked(host string) bool {
	return s.blacklist[host]
}
