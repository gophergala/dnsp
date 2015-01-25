package dnsp

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
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
type HostsReaderFunc func(string)

func (h *HostsReader) ReadFunc(fn HostsReaderFunc) {
	scanner := bufio.NewScanner(h.Reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.SplitN(line, "#", 2)[0]

		parts := strings.Fields(line)
		switch len(parts) {
		case 0: // empty line
			continue
		case 1: // single hostname
			fn(parts[0])
		default: // hosts file like syntax
			if parts[0] == "127.0.0.1" || parts[0] == "::1" {
				for _, host := range parts[1:] {
					fn(host)
				}
			}
		}
	}
}

func readHosts(path string, fn HostsReaderFunc) error {
	if u, err := url.Parse(path); err == nil && u.Host != "" {
		return readHostsURL(path, fn)
	}
	return readHostsFile(path, fn)
}

func readHostsFile(path string, fn HostsReaderFunc) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	NewHostsReader(file).ReadFunc(fn)
	file.Close()
	return nil
}

func readHostsURL(url string, fn HostsReaderFunc) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	NewHostsReader(res.Body).ReadFunc(fn)
	res.Body.Close()
	return nil
}
