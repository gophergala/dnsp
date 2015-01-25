package main

import (
	"strings"

	"github.com/miekg/dns"
)

func init() {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil || len(conf.Servers) == 0 {
		return
	}
	DefaultResolve = strings.Join(conf.Servers, ",")
}
