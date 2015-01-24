package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gophergala/dnsp"
)

var version = "0.0"

func main() {
	o := dnsp.Options{}
	flag.StringVar(&o.Bind, "bind", ":53", "address to bind to")
	flag.Parse()

	s, err := dnsp.NewServer(o)
	if err != nil {
		log.Fatal(err)
	}

	stopServer := func() {
		if err := s.Stop(); err != nil {
			log.Fatal(err)
		}
	}
	defer stopServer() // in case of normal exit

	catch(func(s os.Signal) int {
		log.Printf("Stopping DNS proxy…")
		stopServer()
		return 0
	}, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting DNS proxy on %s…", s.Addr())
	if err = s.Start(); err != nil {
		log.Fatal(err)
	}
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
