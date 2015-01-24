package dnsp_test

import (
	"testing"

	"github.com/gophergala/dnsp"
)

func TestParseHostLine(t *testing.T) {
	t.Parallel()

	var hostLineTests = []struct {
		line string
		ip   string
		host string
	}{
		{"127.0.0.1 ---.chine-li.info", "127.0.0.1", "---.chine-li.info"},
		{"127.0.0.1  ---.chine-li.info  ", "127.0.0.1", "---.chine-li.info"},
		{"127.0.0.1\t\t---.chine-li.info", "127.0.0.1", "---.chine-li.info"},
		{"127.0.0.1 localhost # IPv4", "127.0.0.1", "localhost"},
		{"# Comment line", "", ""},
		{" # Comment with a space at the beginning", "", ""},
		{"   # Comment with spaces at the beginning", "", ""},
		{"   #   ", "", ""},
		{"#", "", ""},
		{"", "", ""},
	}

	for _, test := range hostLineTests {
		result := dnsp.ParseHostLine(test.line)

		if result.IP != test.ip {
			t.Errorf("For test of line %q, IP was %q, expected it to be %q", test.line, result.IP, test.ip)
		}

		if result.Host != test.host {
			t.Errorf("For test of line %q, Host was %q, expected it to be %q", test.line, result.Host, test.host)
		}
	}
}
