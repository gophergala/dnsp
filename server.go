package dnsp

import (
	"log"
	"net"
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
	hostsRX hostsRX

	privateHosts   map[string]struct{}
	privateHostsRX map[string]*regexp.Regexp

	// Information about the hosts file, used for polling:
	hostsFile struct {
		size  int64
		path  string
		mtime time.Time
	}
}

const (
	notIPQuery = 0
	_IP4Query  = 4
	_IP6Query  = 6
)

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
		white:   o.Whitelist != "",
		hosts:   hosts{},
		hostsRX: hostsRX{},

		privateHosts:   map[string]struct{}{},
		privateHostsRX: map[string]*regexp.Regexp{},
	}

	IPv4 := net.ParseIP(o.BlockedIP).To4()
	IPv16 := net.ParseIP(o.BlockedIP).To16()

	hostListPath := o.Whitelist
	if hostListPath == "" {
		hostListPath = o.Blacklist
	}
	s.hostsFile.path = hostListPath
	if err := s.loadHostEntries(); err != nil {
		return nil, err
	}
	if o.Poll != 0 {
		go s.monitorHostEntries(o.Poll)
	}
	s.s.Handler = dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		// If no upstream proxy is present, drop the query:
		if len(o.Resolve) == 0 {
			dns.HandleFailed(w, r)
			return
		}

		if len(r.Question) > 0 {
			q := r.Question[0]

			// Filter Questions:
			if r.Question = s.filter(r.Question); len(r.Question) == 0 {
				IPQuery := s.isIPQuery(q)
				m := new(dns.Msg)
				m.SetReply(r)

				switch IPQuery {
				case _IP4Query:
					rr_header := dns.RR_Header{
						Name:   q.Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    600,
					}

					a := &dns.A{Hdr: rr_header, A: IPv4}
					m.Answer = append(m.Answer, a)

				case _IP6Query:
					rr_header := dns.RR_Header{
						Name:   q.Name,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    600,
					}
					aaaa := &dns.AAAA{Hdr: rr_header, AAAA: IPv16}
					m.Answer = append(m.Answer, aaaa)
				}

				log.Printf("blocked: %s", q.Name)
				w.WriteMsg(m)
				return
			}
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

func (h *Server) isIPQuery(q dns.Question) int {
	if q.Qclass != dns.ClassINET {
		return notIPQuery
	}

	switch q.Qtype {
	case dns.TypeA:
		return _IP4Query
	case dns.TypeAAAA:
		return _IP6Query
	default:
		return notIPQuery
	}
}

// ListenAndServe runs the server
func (s *Server) ListenAndServe() error {
	return s.s.ListenAndServe()
}

// Shutdown stops the server, closing its connection.
func (s *Server) Shutdown() error {
	return s.s.Shutdown()
}
