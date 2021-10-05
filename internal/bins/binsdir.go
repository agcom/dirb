package bins

import (
	"bs/internal/logs"
	"errors"
	"go.uber.org/multierr"
	"io"
	"os"
	"path/filepath"
)

type binsDir string

func New(path string) Repo {
	d := binsDir(filepath.Clean(path))
	return &d
}

func (d *binsDir) New(name string, b io.Reader) error {
	path := d.binPath(name)

	// Create the temp file; acts as a lock and temporary location for incomplete bytes.
	tmpPath := d.tmpPath(name)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrBusy
		} else {
			return multierr.Append(errors.New("open tmp file failed"), err)
		}
	}

	defer func() {
		err = tmpFile.Close()
		if err != nil && !errors.Is(err, os.ErrClosed) {
			logs.Warnf("close tmp file failed; %v", err)
		}

		err = removeForce(tmpPath)
		if err != nil {
			logs.Warnf("remove tmp file failed; %v", err)
		}
	}()

	err = errIfExists(path)
	if err != nil {
		return err
	}

	_, err = tmpFile.ReadFrom(b)
	if err != nil {
		return multierr.Append(errors.New("read/write failed"), err)
	}

	err = tmpFile.Close()
	if err != nil {
		logs.Warnf("close tmp file failed; %v", err)
	}

	err = os.Rename(tmpPath, path)
	if err != nil {
		return multierr.Append(errors.New("rename failed"), err)
	}

	return nil
}

func (d *binsDir) Open(name string) (io.ReadCloser, error) {
	path := d.binPath(name)
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotExist
		} else {
			return nil, multierr.Combine(errors.New("open file failed"), err)
		}
	}

	return f, nil
}

func (d *binsDir) Up(name string, b io.Reader) error {
	path := d.binPath(name)

	// Create the temp file; acts as a lock and temporary location for incomplete bytes.
	tmpPath := d.tmpPath(name)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrBusy
		} else {
			return multierr.Append(errors.New("open tmp file failed"), err)
		}
	}

	defer func() {
		err = tmpFile.Close()
		if err != nil && !errors.Is(err, os.ErrClosed) {
			logs.Warnf("close tmp file failed; %v", err)
		}

		err = removeForce(tmpPath)
		if err != nil {
			logs.Warnf("remove tmp file failed; %v", err)
		}
	}()

	err = errIfNotExist(path)
	if err != nil {
		return err
	}

	_, err = tmpFile.ReadFrom(b)
	if err != nil {
		return multierr.Append(errors.New("read/write failed"), err)
	}

	err = tmpFile.Close()
	if err != nil {
		logs.Warnf("close tmp file failed; %v", err)
	}

	err = os.Rename(tmpPath, path)
	if err != nil {
		return multierr.Append(errors.New("rename failed"), err)
	}

	return nil
}

func (d *binsDir) Rm(name string) error {
	// Create the temp file; acts as a lock.
	tmpPath := d.tmpPath(name)
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_RDONLY, 0660)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ErrBusy
		} else {
			return multierr.Append(errors.New("open tmp file failed"), err)
		}
	}

	defer func() {
		err = tmpFile.Close()
		if err != nil && !errors.Is(err, os.ErrClosed) {
			logs.Warnf("close tmp file failed; %v", err)
		}

		err = removeForce(tmpPath)
		if err != nil {
			logs.Warnf("remove tmp file failed; %v", err)
		}
	}()

	err = os.Remove(d.binPath(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotExist
		} else {
			return multierr.Combine(errors.New("remove bin file failed"), err)
		}
	}

	return nil
}

func (d *binsDir) All() ([]string, error) {
	f, err := os.Open(string(*d))
	if err != nil {
		return nil, multierr.Append(errors.New("open binsDir failed"), err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logs.Warnf("close binsDir failed; %v", err)
		}
	}(f)

	es, err := f.ReadDir(-1)
	if err != nil {
		return nil, multierr.Append(errors.New("read binsDir entries failed"), err)
	}

	ns := make([]string, len(es))[0:]
	for _, e := range es {
		n := e.Name()
		if !e.Type().IsRegular() {
			logs.Warnf("binsDir entry %s is not a regular file", n)
		} else {
			ns = append(ns, n)
		}
	}

	return ns, nil
}

func (d *binsDir) binPath(name string) string {
	return filepath.Join(string(*d), name)
}

func (d *binsDir) tmpPath(name string) string {
	return filepath.Join(string(*d), "."+name+".tmp")
}

func removeForce(path string) error {
	err := os.Remove(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (d *binsDir) tempPathCreate(name string) (string, error) {
	f, err := os.CreateTemp(string(*d), "."+name+".*.tmp")
	if err != nil {
		return "", multierr.Combine(errors.New("create temp file failed"), err)
	}

	err = f.Chmod(0660)
	if err != nil {
		return "", multierr.Combine(errors.New("chmod temp file failed"), err)
	}

	tempPath := f.Name()

	err = f.Close()
	if err != nil {
		logs.Warnf("close temp file failed; %v", err)
	}

	return tempPath, nil
}

func errIfExists(path string) error {
	ex, err := exists(path)
	if err != nil {
		return multierr.Append(errors.New("check file existence failed"), err)
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
		return multierr.Append(errors.New("check file existence failed"), err)
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
