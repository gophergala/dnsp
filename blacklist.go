package dnsp

import "github.com/miekg/dns"

const (
	Unknown block = iota
	Allowed       // whitelisted
	Blocked       // blacklisted
)

type block uint8

// Allow whitelists a hosts.
func (s *Server) Allow(host string) {
	setHost(s.hosts, host, Allowed)
}

// Block blacklists a host.
func (s *Server) Block(host string) {
	setHost(s.hosts, host, Blocked)
}

func setHost(hosts map[string]block, host string, b block) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}
	hosts[host] = b
}

// IsBlocked returns whether a hostname is blocked.
func (s *Server) IsAllowed(host string) bool {
	b := s.hosts[host]
	return s.white && b == Allowed || b != Blocked
}

func (s *Server) filter(qs []dns.Question) []dns.Question {
	result := []dns.Question{}
	for _, q := range qs {
		if s.IsAllowed(q.Name) {
			result = append(result, q)
		}
	}
	return result
}
