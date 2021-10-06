package bins

import (
	"errors"
	"io"
	"path/filepath"
)

var ErrNotExist = errors.New("the bin does not exist")
var ErrExists = errors.New("the bin already exists")
var ErrBusy = errors.New("the bin is already undergoing a change")

// Repo ; bins repository.
// Each bin is identified by a unique name (a string).
// Functions are concurrent safe unless stated otherwise by an implementation.
type Repo interface {
	// New bin; create a new bin.
	// Returns ErrExists if the bin already exists, or ErrBusy if it's undergoing a change.
	New(name string, b io.Reader) error
	// Open the bin with the given name.
	// Might return nil, ErrNotExist.
	Open(name string) (io.ReadCloser, error)
	// Ow ; overwrite the bin with the given name.
	// Might return ErrNotExist, or ErrBusy.
	Ow(name string, b io.Reader) error
	// Rm ; remove the bin with the given name.
	// Might return ErrNotExist, or ErrBusy.
	Rm(name string) error
	// All ; return all available bins' names.
	// Returns nil, error if anything went wrong.
	All() ([]string, error)
}

func NewDirBins(d string) Repo {
	db := Dir(filepath.Clean(d))
	return &db
}
