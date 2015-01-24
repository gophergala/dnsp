package dnsp

import (
	"fmt"
	"strings"
)

type HostEntry struct {
	IP   string
	Host string
}

func (h *HostEntry) String() string {
	return fmt.Sprintf("%v @ %v", h.Host, h.IP)
}

func ParseHostLine(line string) HostEntry {
	result := HostEntry{}

	if len(line) > 0 {
		parts := strings.Fields(line)

		// TODO: More validation might be smart
		if parts[0] != "#" && len(parts) >= 2 {
			result.IP = parts[0]
			result.Host = parts[1]
		}
	}

	return result
}
