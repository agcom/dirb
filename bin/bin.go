// Package bin exports functionality for working with files (CRUD) in a concurrent-safe fashion.
// Bin is an acronym for binary; don't struggle with file operations; focus on the binary.
// Note that the concurrent-safe guarantee is only valid when working with a filesystem that supports atomicity for common file operations (e.g. create, rename, remove, etc.).
package bin

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"io"
	"os"
	"path/filepath"
)

func New(path string, b io.Reader) error {
	return NewLckPath(path, b, DefLckPath(path))
}

func NewLckPath(path string, b io.Reader, lckPath string) error {
	return newOrOverLckPath(true, path, b, lckPath)
}

func NewBare(path string, b io.Reader) error {
	return newOrOverBare(true, path, b)
}

func Open(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, NewErrNotExist(path)
		} else {
			return nil, fmt.Errorf("failed to open %q; %w", path, err)
		}
	}

	return f, nil
}

func Over(path string, b io.Reader) error {
	return OverLckPath(path, b, DefLckPath(path))
}

func OverLckPath(path string, b io.Reader, lckPath string) error {
	return newOrOverLckPath(false, path, b, lckPath)
}

func OverBare(path string, b io.Reader) error {
	return newOrOverBare(false, path, b)
}

func Rm(path string) (rErr error) {
	return RmLckPath(path, DefLckPath(path))
}

func RmLckPath(path string, lckPath string) (rErr error) {
	// Early existence check (not vital)
	err := ErrIfNotExist(path)
	if err != nil {
		return err
	}

	// Rm also need to acquire lock file, to avoid clashing with an ongoing overwrite (Over function).
	lckFile, err := Lck(lckPath)
	if err != nil {
		return err
	}
	defer func() {
		err := Unlck(lckPath, lckFile)
		if err != nil {
			rErr = multierr.Append(rErr, err)
		}
	}()

	return RmBare(path)
}

func RmBare(path string) error {
	err := os.Remove(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewErrNotExist(path)
		} else {
			return fmt.Errorf("failed to remove %q; %w", path, err)
		}
	}

	return nil
}

func newOrOverLckPath(new bool, path string, b io.Reader, lckPath string) (rErr error) {
	// Early existence check (not vital)
	var err error
	if new {
		err = ErrIfExists(path)
	} else {
		err = ErrIfNotExist(path)
	}
	if err != nil {
		return err
	}

	lckFile, err := Lck(lckPath)
	if err != nil {
		return err
	}
	defer func() {
		err := Unlck(lckPath, lckFile)
		if err != nil {
			rErr = multierr.Append(rErr, err)
		}
	}()

	return newOrOverBare(new, path, b)
}

func newOrOverBare(new bool, path string, b io.Reader) (rErr error) {
	var err error
	// Mandatory existence check (if a lock file is involved)
	if new {
		err = ErrIfExists(path)
	} else {
		err = ErrIfNotExist(path)
	}
	if err != nil {
		return err
	}

	// Create and open a temp file (acts as a lock and a temporary place for incomplete bytes).
	// Should reside in the same directory as path to guarantee atomic rename.
	dir, name := filepath.Split(path)
	tmpFile, err := openTmp(dir, "."+name+"-*.tmp")
	tmpPath := tmpFile.Name()
	defer func() {
		tmpFileCloseErr := tmpFile.Close()
		if tmpFileCloseErr != nil {
			if !errors.Is(tmpFileCloseErr, os.ErrClosed) {
				tmpFileCloseErr = fmt.Errorf("failed to close temporary file %q; %w", tmpPath, tmpFileCloseErr)
			} else {
				tmpFileCloseErr = nil
			}
		}

		tmpFileRmErr := os.Remove(tmpPath)
		if tmpFileRmErr != nil {
			if !errors.Is(tmpFileRmErr, os.ErrNotExist) {
				tmpFileRmErr = fmt.Errorf("failed to remove temporary file %q; %w", tmpPath, tmpFileRmErr)
			} else {
				tmpFileRmErr = nil
			}
		}

		rErr = multierr.Combine(rErr, tmpFileRmErr, tmpFileCloseErr)
	}()

	_, err = io.Copy(tmpFile, b)
	if err != nil {
		return fmt.Errorf("failed to read from the given source, or write to temporary file %q; %w", tmpPath, err)
	}

	err = tmpFile.Close()
	var errTmpFileClose error
	if err != nil {
		errTmpFileClose = fmt.Errorf("failed to close temporary file %q; %w", tmpPath, err)
	}

	err = os.Rename(tmpPath, path)
	if err != nil {
		return multierr.Append(fmt.Errorf("failed to rename (move) temporary file %q to %q; %w", tmpPath, path, err), errTmpFileClose)
	}

	return errTmpFileClose
}

func openTmp(dir string, pattern string) (*os.File, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to open a temporary file in directory %q; %w", dir, err)
	}

	err = f.Chmod(0660)
	if err != nil {
		return f, fmt.Errorf("failed to change permission bits of temporary file %q; %w", f.Name(), err)
	}

	return f, nil
}

func DefLckPath(path string) string {
	dir, name := filepath.Split(path)
	return filepath.Join(dir, fmt.Sprintf(".%s.lck.tmp", name))
}
