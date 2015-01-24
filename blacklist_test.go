package dnsp_test

import (
	"testing"

	"github.com/gophergala/dnsp"
	"github.com/miekg/dns"
)

func TestIsBlocked(t *testing.T) {
	t.Parallel()

	s, err := dnsp.NewServer(dnsp.Options{
		Bind: ":0",
	})
	if err != nil {
		t.Fatal(err)
	}

	for host, blocked := range map[string]bool{
		"example.com.": false,
		"blocked.net.": true,
	} {
		if act := s.IsHostBlocked(host); blocked != act {
			t.Errorf("expected s.IsHostBlocked(%q) to be %v, got %v", host, blocked, act)
		}

		m := &dns.Msg{}
		m.SetQuestion(dns.Fqdn(host), dns.TypeA)
		if act := s.IsBlocked(m); blocked != act[0] {
			t.Errorf("expected s.IsBlocked(%T{%q}) to be %v, got %v", m, host, []bool{blocked}, act)
		}
	}
}
