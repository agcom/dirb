package jsns

import (
	"errors"
	"github.com/agcom/bs/bins"
)

var ErrNotExist = errors.New("the json does not exist")
var ErrExists = errors.New("the json already exists")

type Jsn = interface{}
type JsnObj = map[string]Jsn

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

var ErrNotExistObj = errors.New("the json object does not exist")
var ErrExistsObj = errors.New("the json object already exists")

type RepoObj interface {
	New(name string, jo JsnObj) error
	Get(name string) (JsnObj, error)
	Ow(name string, jo JsnObj) error
	Rm(name string) error
	All() ([]string, error)
}

func NewObjJsns(j Repo) *ObjJsns {
	return &ObjJsns{j}
}
