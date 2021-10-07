package bin

import (
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"os"
)

func Lck(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, NewErrLcked(path)
		} else {
			return nil, fmt.Errorf("failed to open or create %q to use as a lock file; %w", path, err)
		}
	}

	return f, nil
}

func Unlck(path string, f *os.File) (rErr error) {
	err := f.Close()
	if err != nil && !errors.Is(err, os.ErrClosed) {
		rErr = fmt.Errorf("failed to close lock file %q; %w", path, err)
	}

	err = os.Remove(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		rErr = multierr.Append(fmt.Errorf("failed to remove lock file %q; %w", path, err), rErr)
	}

	return
}

type ErrLcked string

func (e *ErrLcked) Error() string {
	return fmt.Sprintf("File %q is already locked", e.Path())
}

func (e *ErrLcked) Path() string {
	return string(*e)
}

func NewErrLcked(path string) *ErrLcked {
	e := ErrLcked(path)
	return &e
}
