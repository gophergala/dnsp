package dnsp

import "io"

func ReadConfig(src io.Reader, fn func(string)) {
	readConfig(src, fn)
}