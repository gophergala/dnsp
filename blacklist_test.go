package dnsp_test

import (
	"testing"

	"github.com/gophergala/dnsp"
)

func TestIsAllowed(t *testing.T) {
	t.Parallel()

	s := dnsp.NewServer(dnsp.Options{
		White: true,
	})
	s.Whitelist("google.com")
	s.Whitelist("github.com")

	for host, ok := range map[string]bool{
		"blocked.net.": false,
		"example.com.": false,
		"github.com.":  true,
		"google.com.":  true,
	} {
		if act := s.IsAllowed(host); ok != act {
			t.Errorf("expected s.IsAllowed(%q) to be %v, got %v", host, ok, act)
		}
	}
}
