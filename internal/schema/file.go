package schema

import "io"

type InFile struct {
	Raw  io.Reader
	Name string
	Size int
}

type OutFile struct {
	Raw io.Reader
}
