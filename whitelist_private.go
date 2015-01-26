package dnsp

import "strings"

func (s *Server) privateHostEntries() []string {
	s.m.RLock()
	defer s.m.RUnlock()

	result := []string{}
	for host := range s.privateHosts {
		result = append(result, host)
	}
	for host := range s.privateHostsRX {
		result = append(result, host)
	}
	return result
}

func (s *Server) addPrivateHostEntry(host string) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}

	if !strings.ContainsRune(host, '*') {
		// Plain host string:
		s.m.Lock()
		s.privateHosts[host] = struct{}{}
		s.m.Unlock()
	} else if rx := compilePattern(host); rx != nil {
		// Host pattern (regex):
		s.m.Lock()
		s.privateHostsRX[rx.String()] = rx
		s.m.Unlock()
	}
}

func (s *Server) removePrivateHostEntry(host string) {
	if host == "" {
		return
	}
	if host[len(host)-1] != '.' {
		host += "."
	}

	if !strings.ContainsRune(host, '*') {
		// Plain host string:
		s.m.Lock()
		delete(s.privateHosts, host)
		s.m.Unlock()
	} else if rx := compilePattern(host); rx != nil {
		// Host pattern (regex):
		s.m.Lock()
		delete(s.privateHostsRX, rx.String())
		s.m.Unlock()
	}
}
