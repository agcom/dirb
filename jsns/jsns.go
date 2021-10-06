package jsns

import (
	"errors"
	"github.com/agcom/bs/bins"
)

var ErrNotExist = errors.New("the jsn does not exist")
var ErrExists = errors.New("the jsn already exists")
var ErrBusy = errors.New("the jsn is already undergoing a change")

type Jsn = interface{}

type Repo interface {
	New(name string, j Jsn) error
	Get(name string) (Jsn, error)
	Ow(name string, j Jsn) error
	Rm(name string) error
	All() ([]string, error)
}

func NewBins(b bins.Repo) *Bins {
	return &Bins{b}
}
