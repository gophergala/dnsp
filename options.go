package dnsp

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

// Options can be passed to NewServer().
type Options struct {
	Net     string
	Bind    string
	Resolve []string
	Poll    time.Duration

	Whitelist string
	Blacklist string
}

// validate verifies that the options are correct.
func (o *Options) validate() error {
	if o.Net == "" {
		o.Net = "udp"
	}
	if o.Net != "udp" && o.Net != "tcp" {
		return fmt.Errorf("net: must be one of ‘tcp’, ‘udp’")
	}

	if !strings.Contains(o.Bind, ":") {
		o.Bind += ":53"
	}
	if l := len(o.Bind); l >= 4 && o.Bind[l-4:] == ":dns" {
		o.Bind = o.Bind[:l-4] + ":53"
	}
	if o.Bind[0] == ':' {
		o.Bind = "0.0.0.0" + o.Bind
	}

	for i, res := range o.Resolve {
		if !strings.Contains(res, ":") {
			res += ":53"
		}
		addr, err := net.ResolveUDPAddr("udp", res)
		if err != nil {
			return err
		}
		o.Resolve[i] = addr.String()
	}

	if o.Poll != 0 && o.Poll < time.Second {
		return errors.New("--poll cannot be shorter than 1s")
	}
	o.Poll -= o.Poll % time.Second

	var err error
	if o.Whitelist != "" {
		if o.Blacklist != "" {
			return errors.New("--whitelist and --blacklist are mutually exclusive")
		}
		if o.Whitelist, err = pathOrURL(o.Whitelist); err != nil {
			return err
		}
	}
	if o.Blacklist != "" {
		if o.Blacklist, err = pathOrURL(o.Blacklist); err != nil {
			return err
		}
	}

	return nil
}

func pathOrURL(path string) (string, error) {
	if u, err := url.Parse(path); err == nil && u.Host != "" {
		return u.String(), nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}
	return path, nil
}
