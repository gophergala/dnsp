package dnsp

import (
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// Server implements a DNS server.
type Server struct {
	c *dns.Client
	s *dns.Server

	// White, when set to true, causes the server to work in white-listing mode.
	// It will only resolve queries that have been white-listed.
	//
	// When set to false, it will resolve anything that is not blacklisted.
	white bool

	// Protect access to the hosts file with a mutex.
	m sync.RWMutex

	// A combined whitelist/blacklist. It contains both whitelist and blacklist entries.
	hosts hosts

	// Regex based whitelist and blacklist, depending on the value of `white`.
	hostsRX []*regexp.Regexp
}

// NewServer creates a new Server with the given options.
func NewServer(o Options) (*Server, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	s := Server{
		c: &dns.Client{},
		s: &dns.Server{
			Net:  o.Net,
			Addr: o.Bind,
		},
		white: o.Whitelist != "",
		hosts: hosts{},
	}
	if o.Whitelist != "" {
		if err := s.loadWhitelist(o.Whitelist); err != nil {
			return nil, err
		}
		if o.Poll != 0 {
			go func() {
				for _ = range time.Tick(o.Poll) {
					log.Printf("dnsp: checking %q for updates…", o.Whitelist)
				}
			}()
		}
	}
	if o.Blacklist != "" {
		if err := s.loadBlacklist(o.Blacklist); err != nil {
			return nil, err
		}
		if o.Poll != 0 {
			go func() {
				for _ = range time.Tick(o.Poll) {
					log.Printf("dnsp: checking %q for updates…", o.Blacklist)
				}
			}()
		}
	}
	s.s.Handler = dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		// If no upstream proxy is present, drop the query:
		if len(o.Resolve) == 0 {
			dns.HandleFailed(w, r)
			return
		}

		// Filter Questions:
		if r.Question = s.filter(r.Question); len(r.Question) == 0 {
			w.WriteMsg(r)
			return
		}

		// Proxy Query:
		for _, addr := range o.Resolve {
			in, _, err := s.c.Exchange(r, addr)
			if err != nil {
				continue
			}
			w.WriteMsg(in)
			return
		}
		dns.HandleFailed(w, r)
	})
	return &s, nil
}

// ListenAndServe runs the server
func (s *Server) ListenAndServe() error {
	return s.s.ListenAndServe()
}

// Shutdown stops the server, closing its connection.
func (s *Server) Shutdown() error {
	return s.s.Shutdown()
}
