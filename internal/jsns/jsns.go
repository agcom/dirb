package jsns

import (
	"bs/bins"
	"errors"
)

var ErrNotExist = errors.New("a jsn with the given name doesn't exist")
var ErrExists = errors.New("a jsn with the given name already exists")
var ErrBusy = errors.New("the jsn is undergoing changes")

type Jsn = interface{}

type Repo interface {
	New(name string, j Jsn) error
	Get(name string) (Jsn, error)
	Up(name string, j Jsn) error
	Rm(name string) error
	All() ([]string, error)
}

func NewJsnsBins(b bins.Repo) Repo {
	return &jsnsBins{b}
}
