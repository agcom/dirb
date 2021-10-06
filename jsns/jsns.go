package jsns

import (
	"errors"
)

var ErrNotExist = errors.New("the jsn does not exist")
var ErrExists = errors.New("the jsn already exists")

type Jsn = interface{}

type Repo interface {
	New(name string, j Jsn) error
	Get(name string) (Jsn, error)
	Ow(name string, j Jsn) error
	Rm(name string) error
	All() ([]string, error)
}
