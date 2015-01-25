package dnsp_test

import (
	"log"

	"github.com/gophergala/dnsp"
)

func Example() {
	// Create a server that listens on :1053, on all interfaces.
	// DNS queries will be proxied to Google's public nameservers.
	s, err := dnsp.NewServer(dnsp.Options{
		Bind:    ":1053",
		Resolve: []string{"8.8.4.4", "8.8.8.8"},
		// Block hosts listed in a community-managed file:
		Blacklist: "http://hosts-file.net/download/hosts.txt",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Remember to close it:
	defer s.Shutdown()

	// Start accepting DNS queries:
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
