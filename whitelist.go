package dnsp

import (
	"log"
	"regexp"
	"strings"

	"github.com/miekg/dns"
)

type hosts map[string]struct{}

// isAllowed returns whether we are allowed to resolve this host.
//
// If the server is whitelisting, the result will be true if the host is on the whitelist.
// If the server is blacklisting, the result will be true if the host is NOT on the blacklist.
//
// NOTE: "host" must end with a dot.
func (s *Server) isAllowed(host string) bool {
	s.m.RLock()
	defer s.m.RUnlock()

	_, ok := s.hosts[host]

	if s.white { // whitelist mode
		if ok {
			return true
		}
		for _, rx := range s.hostsRX {
			if rx.MatchString(host) {
				return true
			}
		}
		return false
	}

	// blacklist mode
	if ok {
		return false
	}
	for _, rx := range s.hostsRX {
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

func (s *Server) loadHostEntries(path string) error {
	return readHosts(path, s.addHostEntry)
}

func (s *Server) addHostEntry(host string) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}

	// Plain host string:
	if !strings.ContainsRune(host, '*') {
		s.m.Lock()
		s.hosts[host] = struct{}{}
		s.m.Unlock()
	}

	// Host pattern (regex):
	if pat := compilePattern(host); pat != nil {
		s.m.Lock()
		s.hostsRX = append(s.hostsRX, compilePattern(host))
		s.m.Unlock()
	}
}

func compilePattern(pat string) *regexp.Regexp {
	pat = strings.Replace(pat, ".", `\.`, -1)
	pat = strings.Replace(pat, "*", ".*", -1)
	pat = "^" + pat + `$`
	rx, err := regexp.Compile(pat)
	if err != nil {
		log.Printf("dnsp: could not compile %q: %s", pat, err)
		return nil
	}
	return rx
}
