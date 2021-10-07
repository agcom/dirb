package bin

import (
	"errors"
	"fmt"
	"os"
)

type ErrExists string

func (e *ErrExists) Error() string {
	return fmt.Sprintf("%q already exists", string(*e))
}

func NewErrExists(path string) *ErrExists {
	err := ErrExists(path)
	return &err
}

type ErrNotExist string

func (e *ErrNotExist) Error() string {
	return fmt.Sprintf("%q doesn't exist", string(*e))
}

func NewErrNotExist(path string) *ErrNotExist {
	err := ErrNotExist(path)
	return &err
}

func ErrIfExists(path string) error {
	ex, err := exists(path)
	if err != nil {
		return fmt.Errorf("failed to check if %q exists or not; %w", path, err)
	}

	if ex {
		return NewErrExists(path)
	} else {
		return nil
	}
}

func ErrIfNotExist(path string) error {
	ex, err := exists(path)
	if err != nil {
		return fmt.Errorf("failed to check if %q exists or not; %w", path, err)
	}

	if !ex {
		return NewErrNotExist(path)
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
