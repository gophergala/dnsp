// Package dnsp contains a simple DNS proxy.
package dnsp

import (
	"log"
	"strings"

	"github.com/miekg/dns"
)

// Server implements a DNS server.
type Server struct {
	c *dns.Client
	s *dns.Server

	upstream  string
	blacklist map[string]bool
}

// NewServer creates a new Server with the given options.
func NewServer(o Options) *Server {
	if o.Server != "" && !strings.Contains(o.Server, ":") {
		o.Server += ":53"
	}
	s := Server{
		c: &dns.Client{},
		s: &dns.Server{
			Net:  "udp",
			Addr: o.Bind,
		},
		upstream:  o.Server,
		blacklist: map[string]bool{},
	}
	s.s.Handler = dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		// If no upstream proxy is present, drop the query:
		if o.Server == "" {
			dns.HandleFailed(w, r)
			return
		}

		// Filter Questions:
		if r.Question = s.filter(r.Question); len(r.Question) == 0 {
			w.WriteMsg(r)
			return
		}

		// Proxy Query:
		in, rtt, err := s.c.Exchange(r, o.Server)
		if err != nil {
			log.Printf("error=exchange_failed details=%q", err)
			dns.HandleFailed(w, r)
			return
		}
		log.Printf("debug=exchange_ok rtt=%s", rtt)
		w.WriteMsg(in)

	})
	return &s
}

// ListenAndServe runs the server
func (s *Server) ListenAndServe() error {
	return s.s.ListenAndServe()
}

// Shutdown stops the server, closing any kernel buffers.
func (s *Server) Shutdown() error {
	return s.s.Shutdown()
}
