package dnsp

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type hosts map[string]struct{}

type hostsRX map[string]regexp.Regexp

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

// Load the host entries into separate structures and swap the existing entries.
func (s *Server) loadHostEntries() error {
	hosts := hosts{}
	hostsRX := hostsRX{}

	s.m.RLock()
	path := s.hostsFile.path
	s.m.RUnlock()

	if err := readHosts(path, func(host string) {
		if host[len(host)-1] != '.' {
			host += "."
		}

		if !strings.ContainsRune(host, '*') {
			// Plain host string:
			hosts[host] = struct{}{}
		} else if rx := compilePattern(host); rx != nil {
			// Host pattern (regex):
			hostsRX[rx.String()] = *rx
		}
	}); err != nil {
		return err
	}

	s.m.Lock()
	s.hosts = hosts
	s.hostsRX = hostsRX
	s.m.Unlock()

	return nil
}

func (s *Server) monitorHostEntries(poll time.Duration) {
	s.m.RLock()
	hf := s.hostsFile
	s.m.RUnlock()

	for _ = range time.Tick(poll) {
		log.Printf("dnsp: checking %q for updates…", hf.path)

		mtime, size, err := hostsFileMetadata(hf.path)
		if err != nil {
			log.Printf("dnsp: %s", err)
			continue
		}

		if hf.mtime.Equal(mtime) && hf.size == size {
			continue // no updates
		}

		if err := s.loadHostEntries(); err != nil {
			log.Printf("dnsp: %s", err)
		}
	}
}

func (s *Server) addHostEntry(host string) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}

	if !strings.ContainsRune(host, '*') {
		// Plain host string:
		s.m.Lock()
		s.hosts[host] = struct{}{}
		s.m.Unlock()
	} else if rx := compilePattern(host); rx != nil {
		// Host pattern (regex):
		s.m.Lock()
		s.hostsRX[rx.String()] = *rx
		s.m.Unlock()
	}
}

func (s *Server) removeHostEntry(host string) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}

	if !strings.ContainsRune(host, '*') {
		// Plain host string:
		s.m.Lock()
		delete(s.hosts, host)
		s.m.Unlock()
	} else if rx := compilePattern(host); rx != nil {
		// Host pattern (regex):
		s.m.Lock()
		delete(s.hostsRX, rx.String())
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
