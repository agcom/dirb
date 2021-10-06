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

type ErrTmpRm struct {
	path  string
	cause error
}

// ErrTmpClose does not mean that the operation failed; as long as the tmp file was removed, the operation succeeded; it's usually a warning sign that something is wrong.
type ErrTmpClose struct {
	path  string
	cause error
}

func NewErrTmpRm(path string, cause error) *ErrTmpRm {
	return &ErrTmpRm{
		path:  path,
		cause: cause,
	}
}

func (e *ErrTmpRm) Error() string {
	return fmt.Sprintf("failed to remove \"%s\"", e.path)
}

func (e *ErrTmpRm) Unwrap() error {
	return e.cause
}

func NewErrTmpClose(path string, cause error) *ErrTmpClose {
	return &ErrTmpClose{
		path:  path,
		cause: cause,
	}
}

func (e *ErrTmpClose) Error() string {
	return fmt.Sprintf("failed to close \"%s\"", e.path)
}

func (e *ErrTmpClose) Unwrap() error {
	return e.cause
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

func (d *Dir) Ow(name string, b io.Reader) error {
	return d.newOrOw(false, name, b)
}

func (d *Dir) Rm(name string) (rErr error) {
	path := d.path(name)

	// Early existence check (not vital)
	err := errIfNotExist(path)
	if err != nil {
		return err
	}

	// Create the temp file; acts as a lock.
	tmpPath := d.tmpPath(name)
	tmpFile, err := d.acquireTmpFilePath(tmpPath)
	if err != nil {
		return err
	}

	tmpFileCloseAttempted := false
	tmpFileRm := false

	defer func() {
		if !tmpFileCloseAttempted {
			err := tmpFile.Close()
			tmpFileCloseAttempted = true
			if err != nil {
				rErr = multierr.Append(rErr, NewErrTmpClose(tmpPath, err))
			}
		}

		if !tmpFileRm {
			err := os.Remove(tmpPath)
			if err != nil {
				rErr = multierr.Append(rErr, NewErrTmpRm(tmpPath, err))
			} else {
				tmpFileRm = true
			}
		}
	}()

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

func (d *Dir) newOrOw(new bool, name string, b io.Reader) (rErr error) {
	path := d.path(name)

	// Early existence check (not vital)
	var err error
	if new {
		err = errIfExists(path)
	} else {
		err = errIfNotExist(path)
	}
	if err != nil {
		return err
	}

	// Create the temp file; acts as a lock and temporary location for incomplete bytes.
	tmpPath := d.tmpPath(name)
	tmpFile, err := d.acquireTmpFilePath(tmpPath)
	if err != nil {
		return err
	}

	tmpFileCloseAttempted := false
	tmpFileRm := false

	defer func() {
		if !tmpFileCloseAttempted {
			err := tmpFile.Close()
			tmpFileCloseAttempted = true
			if err != nil {
				rErr = multierr.Append(rErr, NewErrTmpClose(tmpPath, err))
			}
		}

		if !tmpFileRm {
			err := os.Remove(tmpPath)
			if err != nil {
				rErr = multierr.Combine(rErr, NewErrTmpRm(tmpPath, err))
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

	errTmpCloseCause := tmpFile.Close()
	var errTmpClose *ErrTmpClose
	tmpFileCloseAttempted = true
	if errTmpCloseCause != nil {
		errTmpClose = NewErrTmpClose(tmpPath, errTmpCloseCause)
	}

	err = os.Rename(tmpPath, path)
	if err != nil {
		return multierr.Combine(fmt.Errorf("failed to rename (move) \"%s\" to \"%s\"", tmpPath, path), err, errTmpClose)
	}
	tmpFileRm = true

	return errTmpClose
}

func (d *Dir) acquireTmpFile(name string) (*os.File, error) {
	tmpPath := d.tmpPath(name)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, ErrBusy
		} else {
			return nil, multierr.Append(fmt.Errorf("failed to open or create \"%s\"", tmpPath), err)
		}
	}

	return tmpFile, nil
}

func (d *Dir) acquireTmpFilePath(path string) (*os.File, error) {
	tmpFile, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, ErrBusy
		} else {
			return nil, multierr.Append(fmt.Errorf("failed to open or create \"%s\"", path), err)
		}
	}

	return tmpFile, nil
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
