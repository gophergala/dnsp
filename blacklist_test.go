package dnsp_test

import (
	"testing"

	"github.com/gophergala/dnsp"
)

func TestIsAllowed(t *testing.T) {
	t.Parallel()

	for host, blocked := range map[string]bool{
		"example.com.": false,
		"blocked.net.": true,
	} {
		s := dnsp.NewServer(dnsp.Options{})
		if blocked {
			s.Block(host)
		}

		if act := s.IsAllowed(host); blocked != act {
			t.Errorf("expected s.IsAllowed(%q) to be %v, got %v", host, blocked, act)
		}
	}
}
