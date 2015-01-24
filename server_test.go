package dnsp_test

import (
	"testing"

	"github.com/gophergala/dnsp"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	for bind, ok := range map[string]bool{
		":0":              true,  // random port
		"0.0.0.0:0":       true,  // random port, all addresses
		"127.0.0.1:0":     true,  // random port, loopback address
		"127.0.0.1:65537": false, // invalid port
	} {
		s, err := dnsp.NewServer(dnsp.Options{
			Bind: bind,
		})
		if ok {
			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}
			if s == nil {
				t.Errorf("expected a %T, got nil", s)
				continue
			}
			s.Stop()
		} else {
			if err == nil {
				t.Error("expected an error, got nil")
			}
		}
	}
}
