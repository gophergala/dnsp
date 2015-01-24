package dnsp_test

import (
	"testing"

	"github.com/gophergala/dnsp"
)

func TestIsAllowedWhite(t *testing.T) {
	t.Parallel()

	s := dnsp.NewServer(dnsp.Options{
		White: true,
	})
	s.Whitelist("github.com")
	s.Whitelist("google.com")

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

func TestIsAllowedBlack(t *testing.T) {
	t.Parallel()

	s := dnsp.NewServer(dnsp.Options{
		White: false,
	})
	s.Blacklist("doubleclick.net")
	s.Blacklist("porn.com")

	for host, ok := range map[string]bool{
		"doubleclick.net.": false,
		"example.com.":     true,
		"github.com.":      true,
		"google.com.":      true,
		"porn.com.":        false,
	} {
		if act := s.IsAllowed(host); ok != act {
			t.Errorf("expected s.IsAllowed(%q) to be %v, got %v", host, ok, act)
		}
	}
}
