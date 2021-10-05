package bins

import (
	"errors"
	"io"
)

var ErrNotExist = errors.New("a bin with the given name doesn't exist")
var ErrExists = errors.New("a bin with the given name already exists")
var ErrBusy = errors.New("the bin is undergoing changes")

type Repo interface {
	New(name string, b io.Reader) error
	Open(name string) (io.ReadCloser, error)
	Up(name string, b io.Reader) error
	Rm(name string) error
	All() ([]string, error)
}

func NewBinsDir(d string) Repo {
	return New(d)
}
