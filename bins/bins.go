package bins

import (
	"errors"
	"io"
)

var ErrNotExist = errors.New("the binary does not exist")
var ErrExists = errors.New("the binary already exists")

// Repo ; bins repository.
// Each bin is identified by a unique name (a string).
// Functions are concurrent safe unless stated otherwise by an implementation.
type Repo interface {
	// New bin; create a new bin.
	// Returns ErrExists if the bin already exists.
	New(name string, b io.Reader) error
	// Open the bin with the given name.
	// Might return nil, ErrNotExist.
	Open(name string) (io.ReadCloser, error)
	// Ow ; overwrite the bin with the given name.
	// Might return ErrNotExist.
	Ow(name string, b io.Reader) error
	// Rm ; remove the bin with the given name.
	// Might return ErrNotExist.
	Rm(name string) error
	// All ; return all available bins' names.
	All() ([]string, error)
}
