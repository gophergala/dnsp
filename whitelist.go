package dnsp

import (
	"crypto/md5"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type hosts map[checksum]struct{}

type hostsRX map[checksum]*regexp.Regexp

type checksum [md5.Size / 2]byte

// isAllowed returns whether we are allowed to resolve this host.
//
// If the server is whitelisting, the result will be true if the host is on the whitelist.
// If the server is blacklisting, the result will be true if the host is NOT on the blacklist.
//
// NOTE: "host" must end with a dot.
func (s *Server) isAllowed(host string) bool {
	s.m.RLock()
	defer s.m.RUnlock()

	_, ok := s.hosts[hash(host)]
	if !ok {
		_, ok = s.privateHosts[host]
	}

	if s.white { // whitelist mode
		if ok {
			return true
		}
		for _, rx := range s.hostsRX {
			if rx.MatchString(host) {
				return true
			}
		}
		for _, rx := range s.privateHostsRX {
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
	for _, rx := range s.privateHostsRX {
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
	s.m.RLock()
	path := s.hostsFile.path
	s.m.RUnlock()

	if path == "" {
		return nil
	}

	hosts := hosts{}
	hostsRX := hostsRX{}

	if err := readHosts(path, func(host string) {
		if host[len(host)-1] != '.' {
			host += "."
		}

		if !strings.ContainsRune(host, '*') {
			// Plain host string:
			hosts[hash(host)] = struct{}{}
		} else if rx := compilePattern(host); rx != nil {
			// Host pattern (regex):
			hostsRX[hash(rx.String())] = rx
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

	if hf.path == "" {
		return
	}

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

		s.m.Lock()
		s.hostsFile.mtime = mtime
		s.hostsFile.size = size
		hf = s.hostsFile
		s.m.Unlock()
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
		s.hosts[hash(host)] = struct{}{}
		s.m.Unlock()
	} else if rx := compilePattern(host); rx != nil {
		// Host pattern (regex):
		s.m.Lock()
		s.hostsRX[hash(rx.String())] = rx
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
		delete(s.hosts, hash(host))
		s.m.Unlock()
	} else if rx := compilePattern(host); rx != nil {
		// Host pattern (regex):
		s.m.Lock()
		delete(s.hostsRX, hash(rx.String()))
		s.m.Unlock()
	}
}

func (s *Server) publicEntriesCount() int {
	s.m.Lock()
	n := len(s.hosts) + len(s.hostsRX)
	s.m.Unlock()
	return n
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

func hash(s string) checksum {
	sum := md5.Sum([]byte(s))
	return checksum{
		sum[0] ^ sum[1],
		sum[2] ^ sum[3],
		sum[4] ^ sum[5],
		sum[6] ^ sum[7],
		sum[8] ^ sum[9],
		sum[10] ^ sum[11],
		sum[12] ^ sum[13],
		sum[14] ^ sum[15],
	}
}
