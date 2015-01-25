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

// Whitelist whitelists a host or a pattern.
func (s *Server) Whitelist(host string) {
	if strings.ContainsRune(host, '*') {
		s.rxWhitelist = appendPattern(s.rxWhitelist, host)
	} else {
		setHost(s.hosts, host, white)
	}
}

// Blacklist blacklists a host.
// If rx is true, the given host is treated as a regular expression.
func (s *Server) Blacklist(host string) {
	if strings.ContainsRune(host, '*') {
		s.rxBlacklist = appendPattern(s.rxBlacklist, host)
	} else {
		setHost(s.hosts, host, black)
	}
}

func setHost(hosts map[string]host, host string, b host) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}
	hosts[host] = b
}

// IsAllowed returns whether we are allowed to resolve this host.
//
// If the server is whitelisting, the rusilt will be true if the host is on the whitelist.
// If the server is blacklisting, the result will be true if the host is NOT on the blacklist.
//
// NOTE: "host" must end with a dot.
func (s *Server) IsAllowed(host string) bool {
	b := s.hosts[host]
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
		if s.IsAllowed(q.Name) {
			result = append(result, q)
		}
	}
	return result
}

func (s *Server) loadWhitelist(path string) error {
	return readHosts(path, s.Whitelist)
}

func (s *Server) loadBlacklist(path string) error {
	return readHosts(path, s.Blacklist)
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
