package dnsp

import (
	"log"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

const (
	white = iota + 1 // whitelisted
	black            // blacklisted
)

type host uint8

type hosts map[string]host

// isAllowed returns whether we are allowed to resolve this host.
//
// If the server is whitelisting, the result will be true if the host is on the whitelist.
// If the server is blacklisting, the result will be true if the host is NOT on the blacklist.
//
// NOTE: "host" must end with a dot.
func (s *Server) isAllowed(host string) bool {
	s.m.RLock()
	b := s.hosts[host]
	s.m.RUnlock()
	if s.white { // check whitelists
		if b == white {
			return true
		}
		for _, rx := range s.rxWhitelist {
			if rx.MatchString(host) {
				return true
			}
		}
		return false
	}
	// check blacklists
	if b == black {
		return false
	}
	for _, rx := range s.rxBlacklist {
		if rx.MatchString(host) {
			return false
		}
	}
	return true
}

func (s *Server) filter(qs []dns.Question) []dns.Question {
	result := []dns.Question{}
	for _, q := range qs {
		if s.isAllowed(q.Name) {
			result = append(result, q)
		}
	}
	return result
}

// whitelist whitelists a host or a pattern.
func (s *Server) whitelist(host string) {
	if strings.ContainsRune(host, '*') {
		s.rxWhitelist = appendPattern(s.rxWhitelist, host)
	} else {
		s.setHost(host, white)
	}
}

// blacklist blacklists a host or a pattern.
func (s *Server) blacklist(host string) {
	if strings.ContainsRune(host, '*') {
		s.rxBlacklist = appendPattern(s.rxBlacklist, host)
	} else {
		s.setHost(host, black)
	}
}

func (s *Server) setHost(host string, b host) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}
	s.m.Lock()
	s.hosts[host] = b
	s.m.Unlock()
}

func (s *Server) loadWhitelist(path string) error {
	return readHosts(path, s.whitelist)
}

func (s *Server) loadBlacklist(path string) error {
	return readHosts(path, s.blacklist)
}

func appendPattern(rx []*regexp.Regexp, pat string) []*regexp.Regexp {
	if pat == "" {
		return rx
	}

	pat = strings.Replace(pat, ".", `\.`, -1)
	pat = strings.Replace(pat, "*", ".*", -1)
	pat = "^" + pat + `\.$`
	if r, err := regexp.Compile(pat); err != nil {
		log.Printf("dnsp: could not compile %q: %s", pat, err)
	} else {
		rx = append(rx, r)
	}
	return rx
}
