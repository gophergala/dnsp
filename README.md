# `dnsp`: A DNS Proxy

[![Wercker](https://app.wercker.com/status/f42156d5f4e863ebe8cf0c311bd7800a/s/master "wercker status")](https://app.wercker.com/project/bykey/f42156d5f4e863ebe8cf0c311bd7800a)
[![GoDoc](https://godoc.org/github.com/gophergala/dnsp?status.svg)](https://godoc.org/github.com/gophergala/dnsp)
[![Coverage](http://gocover.io/_badge/github.com/gophergala/dnsp)](http://gocover.io/github.com/gophergala/dnsp)

> `dnsp` is a lightweight but powerful DNS server. Queries are blocked or
> resolved based on a blacklist or a whitelist. Wildcard host patterns are
> supported (e.g. `*.com`) as well as hosted, community-managed hosts files.
> Ideal for running on mobile devices or embedded systems, given its [low
> memory footprint][1] and simple web interface.


### Installation

```sh
$ go get -u github.com/gophergala/dnsp
```

### Example Usage

* Forward all queries to Google's public nameservers:

```sh
$ sudo dnsp --resolve 8.8.4.4,8.8.8.8
```

* Use a community-managed blacklist from [hosts-file.net] and check it hourly
  for changes:

```sh
$ sudo dnsp --blacklist=http://hosts-file.net/download/hosts.txt --poll 1h
```

* Block everything except Wikipedia:

```sh
$ cat > /etc/dnsp.whitelist << EOF
*.wikipedia.org
*.wikimedia.org
wikipedia.org
wikimedia.org
EOF

$ sudo dnsp -r 8.8.8.8 --whitelist=/etc/dnsp.whitelist
```


### Advanced Usage

```sh
$ dnsp -h
NAME:
   dnsp - DNS proxy with whitelist/blacklist support

USAGE:
   dnsp [global options] command [command options] [arguments...]

VERSION:
   0.9.2

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --net, -n "udp"          listen protocol (‘tcp’ or ‘udp’) [$DNSP_NET]
   --listen, -l ":dns"      listen address (host:port, host or :port) [$DNSP_BIND]
   --resolve, -r "8.8.4.4"  comma-separated list of name servers (host:port or host) [$DNSP_SERVER]
   --whitelist, -w          URL or path to file containing whitelisted hosts [$DNSP_WHITELIST]
   --blacklist, -b          URL or path to file containing blacklisted hosts [$DNSP_BLACKLIST]
   --poll, -p "0"           poll the whitelist or blacklist for updates [$DNSP_POLL]
   --http, -t               start a web-based UI on the given address (host:port, host or port) [$DNSP_HTTP]
   --help, -h               show help
   --version, -v            print the version
```

**Notes:**

* `--listen` defaults to `:dns`, which is equivalent to `0.0.0.0:53`, meaning:
  listen on all interfaces, on port 53 (default DNS port). 
* `--resolve` defaults to the list of nameservers found in `/etc/resolv.conf`.
  If no nameservers were found, or the file does not exist (e.g. on Windows),
  the default value will be `8.8.4.4,8.8.8.8" (Google's public DNS service).
  * However, explicitly setting `--resolve` to `false` or an empty string
	disables resolving completely. What that means is all queries will still be
	checked against the active whitelist or blacklist, but ones that would not
	be blocked will return a failure response (as opposed to no response).
* `--whitelist` and `--blacklist` are mutually exclusive. Setting both is an error.
* `--whitelist` and `--blacklist` files are parsed according to a simple syntax:
  * Empty lines are ignored, and `#` begins a single-line comment.
  * Each line can contain a single hostname to be whitelisted or blacklisted.
  * Alternatively, a line can contain a pattern like `*.wikipedia.org` or
	`*.xxx`.
  * Additionally, the `/etc/hosts`-like syntax is supported.
	* However, only lines starting with `127.0.0.1` or `::1` are taken into
	  parsed, everything else is ignored.
	* This is for compatibility with popular, regularly updated blocklists like
	  the ones on [hosts-file.net].
* `--whitelist` and `--blacklist` support both file paths and URLs.
* `--poll` instructs `dnsp` to periodically check the whitelist or blacklist
  file for changes.
  * The file is only re-parsed if the file size or modification time has
    changed since the last read.
  * Same is true for URLs: the `Content-Length` and `Last-Modified` headers are
    compared to previous values before re-downloading the file.


### Running with a non-root user

Because `dnsp` binds to port 53 by default, it requires to be run with a
privileged user on most systems. To avoid having to run `dnsp` with sudo, you
can set the `setuid` and `setgid` access right flags on the compiled
executable:

```
sudo mkdir -p /usr/local/bin
sudo cp $GOPATH/bin/dnsp
sudo chmod ug+s /usr/local/bin/dnsp
```

While `dnsp` will still run with root privileges, at least now we can run it
with a non-admin user (someone who is not in the `sudoers` group).


### But… Why‽

Why, you ask, is a DNS proxy useful?

* It is a simple solution for blocking websites (like [AdBlock]).
* Does not require an HTTP proxy or a SOCKS proxy. Some apps don't like that.
* Easy to set up for mobile devices. Run `dnsmasq` on your router or in any
  embedded Linux system, and configure your home router to use it as the DNS
  server in DHCP responses. The blocklist will now apply to everyone on the
  network.
* Safer than `dnsmasq` for community managed hosts files. Because `dnsp`
  doesn't do any rewriting (it either blocks or proxies), you don't have to
  trust everyone having access to online hosts files not to redirect your
  bank's website to their own servers.

![dnsp](https://cloud.githubusercontent.com/assets/196617/5892473/cc29afe2-a4bf-11e4-9c6a-d1cc5169d62a.png)


[1]: https://github.com/gophergala/dnsp/pull/15#issue-55432972
[hosts-file.net]: http://hosts-file.net
[AdBlock]: https://getadblock.com
