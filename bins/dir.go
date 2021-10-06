package bins

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"io"
	"os"
	"path/filepath"
)

// Dir implements Repo that saves each bin in a file in the specified directory.
// Note that any name argument should be a valid file name (e.g. no filepath.Separator within); otherwise, strange things will happen.
// All functions are concurrent safe, unless the underlying file system doesn't support atomicity for common file operations (e.g. create, rename, remove, etc.).
type Dir string

func (d *Dir) newOrOw(new bool, name string, b io.Reader) (rErr error) {
	path := d.path(name)

	// Create the temp file; acts as a lock and temporary location for incomplete bytes.
	tmpPath := d.tmpPath(name)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrBusy
		} else {
			return multierr.Append(fmt.Errorf("failed to open or create \"%s\"", tmpPath), err)
		}
	}

	tmpFileCloseAttempted := false
	tmpFileRm := false

	defer func() {
		if !tmpFileCloseAttempted {
			err := tmpFile.Close()
			tmpFileCloseAttempted = true
			if err != nil {
				// Ignore the tmp file close error; should be already closed.
			}
		}

		if !tmpFileRm {
			err := os.Remove(tmpPath)
			if err != nil {
				rErr = multierr.Append(fmt.Errorf("failed to remove \"%s\"", tmpPath), err)
			} else {
				tmpFileRm = true
			}
		}
	}()

	if new {
		err = errIfExists(path)
	} else {
		err = errIfNotExist(path)
	}
	if err != nil {
		return err
	}

	_, err = io.Copy(tmpFile, b)
	if err != nil {
		return multierr.Append(fmt.Errorf("failed to read from the given source, or write to \"%s\"", tmpPath), err)
	}

	err = tmpFile.Close()
	tmpFileCloseAttempted = true
	if err != nil {
		// Ignore the tmp file close error; should be already closed.
	}

	err = os.Rename(tmpPath, path)
	if err != nil {
		return multierr.Append(fmt.Errorf("failed to rename (move) \"%s\" to \"%s\"", tmpPath, path), err)
	}
	tmpFileRm = true

	return nil
}

func (d *Dir) New(name string, b io.Reader) error {
	return d.newOrOw(true, name, b)
}

func (d *Dir) Open(name string) (io.ReadCloser, error) {
	path := d.path(name)
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotExist
		} else {
			return nil, multierr.Combine(fmt.Errorf("failed to open \"%s\"", path), err)
		}
	}

	return f, nil
}

func (d *Dir) Ow(name string, b io.Reader) (rErr error) {
	return d.newOrOw(false, name, b)
}

func (d *Dir) Rm(name string) (rErr error) {
	// Create the temp file; acts as a lock.
	tmpPath := d.tmpPath(name)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrBusy
		} else {
			return multierr.Append(fmt.Errorf("failed to open or create \"%s\"", tmpPath), err)
		}
	}

	tmpFileCloseAttempted := false
	tmpFileRm := false

	defer func() {
		if !tmpFileCloseAttempted {
			err := tmpFile.Close()
			tmpFileCloseAttempted = true
			if err != nil {
				// Ignore the tmp file close error; should be already closed.
			}
		}

		if !tmpFileRm {
			err := os.Remove(tmpPath)
			if err != nil {
				rErr = multierr.Append(fmt.Errorf("failed to remove \"%s\"", tmpPath), err)
			} else {
				tmpFileRm = true
			}
		}
	}()

	path := d.path(name)
	err = os.Remove(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotExist
		} else {
			return multierr.Combine(fmt.Errorf("failed to remove \"%s\"", path), err)
		}
	}

	return nil
}

func (d *Dir) All() (rNs []string, rErr error) {
	return d.AllRepIrreg(func(string) {
		// No op; ignore irregular files.
	})
}

// AllRepIrreg is the same as All but reports irregular files through the given function.
func (d *Dir) AllRepIrreg(rep func(string)) (rNs []string, rErr error) {
	dPath := d.Dir()
	f, err := os.Open(dPath)
	if err != nil {
		return nil, multierr.Append(fmt.Errorf("failed to open \"%s\"", dPath), err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			rNs = nil
			rErr = multierr.Append(fmt.Errorf("failed to close \"%s\"", dPath), err)
		}
	}(f)

	es, err := f.ReadDir(-1)
	if err != nil {
		return nil, multierr.Append(fmt.Errorf("failed to read directory entries of \"%s\"", dPath), err)
	}

	ns := make([]string, len(es))[0:]
	for _, e := range es {
		n := e.Name()
		if e.Type().IsRegular() {
			ns = append(ns, n)
		} else {
			rep(n)
		}
	}

	return ns, nil
}

func (d *Dir) Dir() string {
	return string(*d)
}

func (d *Dir) path(name string) string {
	return filepath.Join(d.Dir(), name)
}

func (d *Dir) tmpPath(name string) string {
	return filepath.Join(d.Dir(), "."+name+".tmp")
}

func errIfExists(path string) error {
	ex, err := exists(path)
	if err != nil {
		return multierr.Append(fmt.Errorf("failed to check if \"%s\" exists or not", path), err)
	}

	if ex {
		return ErrExists
	} else {
		return nil
	}
}

func errIfNotExist(path string) error {
	ex, err := exists(path)
	if err != nil {
		return multierr.Append(fmt.Errorf("failed to check if \"%s\" exists or not", path), err)
	}

	if !ex {
		return ErrNotExist
	} else {
		return nil
	}
}

func exists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
