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
	flag.StringVar(&o.Server, "server", "8.8.8.8", "address to proxy to")
	flag.Parse()

	s := dnsp.NewServer(o)

	catch(func(sig os.Signal) int {
		log.Printf("debug=dnsp_shutdown signal=%s", sig)
		s.Shutdown()
		return 0
	}, syscall.SIGINT, syscall.SIGTERM)
	defer s.Shutdown() // in case of normal exit

	log.Printf("debug=dnsp_start bind=%s server=%s", o.Bind, o.Server)
	if err := s.ListenAndServe(); err != nil {
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
