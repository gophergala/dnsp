package dnsp

// Options can be passed to NewServer().
type Options struct {
	// Bind is the local address to listen on.
	Bind string
	// Server is the address of the upstream DNS server to proxy to.
	// If empty, all queries will fail.
	Server string
}
