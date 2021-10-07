package bin

import (
	"fmt"
	"go.uber.org/multierr"
	"io"
	"os"
	"path/filepath"
)

// Dir is just a fancy wrapper around global bin functions; a binary repository that saves all bins in a specified directory.
// Note that any name argument should be a valid file name (e.g. no filepath.Separator within); otherwise, strange things will happen.
type Dir string

func NewDir(d string) *Dir {
	dir := Dir(d)
	return &dir
}

func (d *Dir) New(name string, b io.Reader) error {
	path := d.Path(name)
	return New(path, b)
}

func (d *Dir) Open(name string) (*os.File, error) {
	path := d.Path(name)
	return Open(path)
}

func (d *Dir) Over(name string, b io.Reader) error {
	path := d.Path(name)
	return Over(path, b)
}

func (d *Dir) Rm(name string) (rErr error) {
	path := d.Path(name)
	return Rm(path)
}

func (d *Dir) All() (rNs []string, rErr error) {
	dPath := d.Dir()

	f, err := os.Open(dPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory %q; %w", dPath, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			rErr = multierr.Append(rErr, fmt.Errorf("failed to close directory %q; %w", dPath, err))
		}
	}()

	es, err := f.ReadDir(-1)
	if err != nil {
		rErr = multierr.Append(rErr, fmt.Errorf("failed to read directory entries of %q; %w", dPath, err))
	}

	rNs = make([]string, len(es))[:0]
	for _, e := range es {
		n := e.Name()
		if e.Type().IsRegular() {
			rNs = append(rNs, n)
		} else {
			rErr = multierr.Append(rErr, fmt.Errorf("irregular file %q in the binarys' data directory %q", e.Name(), dPath))
		}
	}

	return
}

func (d *Dir) Dir() string {
	return string(*d)
}

func (d *Dir) Path(name string) string {
	return filepath.Join(d.Dir(), name)
}
