package jsnobjs

import (
	"errors"
	"github.com/agcom/bs/jsns"
)

var ErrNotExist = errors.New("the json object does not exist")
var ErrExists = errors.New("the json object already exists")

type Jsnobj = map[string]jsns.Jsn

type Repo interface {
	New(name string, jo Jsnobj) error
	Get(name string) (Jsnobj, error)
	Ow(name string, jo Jsnobj) error
	Rm(name string) error
	All() ([]string, error)
}
