package dnsp

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

// readHosts reads hosts from a URL or a file.
func readHosts(path string, fn func(string)) error {
	if u, err := url.Parse(path); err == nil && u.Host != "" {
		return readHostsURL(path, fn)
	}
	return readHostsFile(path, fn)
}

// readHostsFile reads hosts from a file.
func readHostsFile(path string, fn func(string)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	readConfig(file, fn)
	file.Close()
	return nil
}

// readHostsURL reads hosts from a URL.
func readHostsURL(url string, fn func(string)) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	readConfig(res.Body, fn)
	res.Body.Close()
	return nil
}

// hostsFileMetadata returns metadata about the hosts file.
//
// If path is a URL, it will make a HEAD request to get the headers (leaving
// the TCP connection open for subsequent GETs, if needed).
//
// If path is a file path, the mtime and file size is returned.
func hostsFileMetadata(path string) (time.Time, int64, error) {
	if u, err := url.Parse(path); err == nil && u.Host != "" {
		res, err := http.Head(u.String())
		if err != nil {
			return time.Time{}, 0, err
		}
		size, _ := strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
		mtime, _ := time.Parse(time.RFC1123, res.Header.Get("Last-Modified"))
		return mtime, size, nil
	}

	fi, err := os.Stat(path)
	if err != nil {
		return time.Time{}, 0, err
	}

	return fi.ModTime(), fi.Size(), nil
}
