package dnsp

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
)

type HostsReader struct {
	Reader io.Reader
}

func NewHostsReader(r io.Reader) *HostsReader {
	return &HostsReader{Reader: r}
}

// HostsReaderFunc is a function that takes a hostname as its first argument.
// The second argument indicates whether the first argument is a regex pattern.
type HostsReaderFunc func(string, bool)

func (h *HostsReader) ReadFunc(fn HostsReaderFunc) {
	scanner := bufio.NewScanner(h.Reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.SplitN(line, "#", 2)[0]

		parts := strings.Fields(line)
		switch len(parts) {
		case 0: // empty line
			continue
		case 1: // hostname or regex
			host := parts[0]
			rx := strings.Contains(host, "*")
			if rx { // prepare regular expression
				host = strings.Replace(host, ".", `\.`, -1)
				host = strings.Replace(host, "*", ".*", -1)
				host = "^" + host + "$"
			}
			fn(host, rx)
		default: // hosts file like syntax
			if parts[0] == "127.0.0.1" || parts[0] == "::1" {
				for _, host := range parts[1:] {
					fn(host, false)
				}
			}
		}
	}
}

func ReadHostFile(path string, fn HostsReaderFunc) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	NewHostsReader(file).ReadFunc(fn)
	file.Close()
	return nil
}

func ReadHostURL(url string, fn HostsReaderFunc) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	NewHostsReader(res.Body).ReadFunc(fn)
	res.Body.Close()
	return nil
}
