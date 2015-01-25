package dnsp

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func readConfig(src io.Reader, fn func(string)) {
	scanner := bufio.NewScanner(src)
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
			if parts[0] == "127.0.0.1" || parts[0] == "0.0.0.0" || parts[0] == "::1" {
				for _, host := range parts[1:] {
					fn(host)
				}
			}
		}
	}
}

func readHosts(path string, fn func(string)) error {
	if u, err := url.Parse(path); err == nil && u.Host != "" {
		return readHostsURL(path, fn)
	}
	return readHostsFile(path, fn)
}

func readHostsFile(path string, fn func(string)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	readConfig(file, fn)
	file.Close()
	return nil
}

func readHostsURL(url string, fn func(string)) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	readConfig(res.Body, fn)
	res.Body.Close()
	return nil
}
