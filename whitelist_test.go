package dnsp_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gophergala/dnsp"
)

func TestIsAllowedWhite(t *testing.T) {
	t.Parallel()

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	tmp.Write([]byte("github.com\n"))
	tmp.Write([]byte("google.com\n"))

	s, err := dnsp.NewServer(dnsp.Options{
		Whitelist: tmp.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}

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

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	tmp.Write([]byte("doubleclick.net\n"))
	tmp.Write([]byte("porn.com\n"))

	s, err := dnsp.NewServer(dnsp.Options{
		Blacklist: tmp.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}

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
