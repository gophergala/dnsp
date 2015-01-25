package dnsp

import "io"

func ReadConfig(src io.Reader, fn func(string)) {
	readConfig(src, fn)
}

func (s *Server) IsAllowed(host string) bool {
	return s.isAllowed(host)
}

func (o *Options) Validate() error {
	return o.validate()
}

func PathOrURL(path string) (string, error) {
	return pathOrURL(path)
}
