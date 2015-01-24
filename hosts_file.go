package dnsp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type HostEntry struct {
	IP   string
	Host string
}

func (h *HostEntry) String() string {
	return fmt.Sprintf("%v @ %v", h.Host, h.IP)
}

func ParseHostLine(line string) *HostEntry {
	result := HostEntry{}

	if len(line) > 0 {
		parts := strings.Fields(line)

		// TODO: More validation might be smart
		if parts[0] != "#" && len(parts) >= 2 {
			result.IP = parts[0]
			result.Host = parts[1]
		}
	}

	return &result
}

type HostsReader struct {
	Reader io.Reader
}

func NewHostsReader(r io.Reader) *HostsReader {
	return &HostsReader{Reader: r}
}

type HostsReaderFunc func(*HostEntry)

func (h *HostsReader) ReadFunc(f HostsReaderFunc) {
	scanner := bufio.NewScanner(h.Reader)
	for scanner.Scan() {
		hostEntry := ParseHostLine(scanner.Text())
		if hostEntry.IP != "" {
			f(hostEntry)
		}
	}
}

func (h *HostsReader) ReadAll() []*HostEntry {
	result := make([]*HostEntry, 0)

	h.ReadFunc(func(h *HostEntry) {
		result = append(result, h)
	})

	return result
}

func ReadHostFile(filename string, f HostsReaderFunc) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	reader := NewHostsReader(file)
	reader.ReadFunc(f)

	return nil
}
