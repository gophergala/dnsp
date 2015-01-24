package dnsp_test

import (
	"strings"
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

func TestHostsReaderReadAll(t *testing.T) {
	t.Parallel()

	hostsFile := strings.NewReader(`127.0.0.1 ---.chine-li.info
127.0.0.1 -ads.avast.dwnldfr.com
# Comment
127.0.0.1 -reports.com-57o.net
127.0.0.1 0-29.com`)

	reader := dnsp.NewHostsReader(hostsFile)
	result := reader.ReadAll()

	if len(result) != 4 {
		t.Errorf("Result is not of expected length 4")
	}

	if result[0].IP != "127.0.0.1" || result[0].Host != "---.chine-li.info" {
		t.Errorf("First entry is wrong")
	}
}
