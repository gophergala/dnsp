package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/gophergala/dnsp"
)

// DefaultResolve is the default list of nameservers for the `--resolve` flag.
var DefaultResolve = "8.8.4.4,8.8.8.8"

func main() {
	app := cli.NewApp()
	app.Name = "dnsp"
	app.Usage = "DNS proxy with whitelist/blacklist support"
	app.Version = "0.0.0"
	app.Author, app.Email = "", ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "net, n",
			Value:  "udp",
			Usage:  "listen protocol (‘tcp’ or ‘udp’)",
			EnvVar: "DNSP_NET",
		},
		cli.StringFlag{
			Name:   "listen, l",
			Value:  ":dns",
			Usage:  "listen address (host:port, host or :port)",
			EnvVar: "DNSP_BIND",
		},
		cli.StringFlag{
			Name:   "resolve, r",
			Value:  DefaultResolve,
			Usage:  "comma-separated list of name servers (host:port or host)",
			EnvVar: "DNSP_SERVER",
		},
		cli.StringFlag{
			Name:   "whitelist, w",
			Usage:  "URL or path to file containing whitelisted hosts",
			EnvVar: "DNSP_WHITELIST",
		},
		cli.StringFlag{
			Name:   "blacklist, b",
			Usage:  "URL or path to file containing blacklisted hosts",
			EnvVar: "DNSP_BLACKLIST",
		},
	}
	app.Action = func(c *cli.Context) {
		resolve := []string{}
		if res := c.String("resolve"); res != "false" && res != "" {
			resolve = strings.Split(res, ",")
		}
		o := &dnsp.Options{
			Net:       c.String("net"),
			Bind:      c.String("listen"),
			Resolve:   resolve,
			Whitelist: c.String("whitelist"),
			Blacklist: c.String("blacklist"),
		}
		if err := o.Validate(); err != nil {
			cli.ShowAppHelp(c)
			log.Fatalf("dnsp: %s", err)
		}
		s, err := dnsp.NewServer(*o)
		if err != nil {
			log.Fatalf("dnsp: %s", err)
		}

		catch(func(sig os.Signal) int {
			os.Stderr.Write([]byte{'\r'})
			log.Printf("dnsp: shutting down")
			s.Shutdown()
			return 0
		}, syscall.SIGINT, syscall.SIGTERM)
		defer s.Shutdown() // in case of normal exit

		if len(o.Resolve) == 0 {
			log.Printf("dnsp: listening on %s", o.Bind)
		} else {
			log.Printf("dnsp: listening on %s, proxying to %s", o.Bind, o.Resolve)
		}
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("dnsp: %s", err)
		}
	}
	app.Run(os.Args)
}

// catch handles system calls using the given handler function.
func catch(handler func(os.Signal) int, signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	for _, s := range signals {
		signal.Notify(c, s)
	}
	go func() {
		os.Exit(handler(<-c))
	}()
}
