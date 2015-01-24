package dnsp_test

import (
	"testing"

	"github.com/gophergala/dnsp"
)

func TestIsBlocked(t *testing.T) {
	t.Parallel()

	s := dnsp.NewServer(dnsp.Options{})

	for host, blocked := range map[string]bool{
		"example.com.": false,
		"blocked.net.": true,
	} {
		if blocked {
			s.Block(host)
		}

		if act := s.IsBlocked(host); blocked != act {
			t.Errorf("expected s.IsBlocked(%q) to be %v, got %v", host, blocked, act)
		}
	}
}
