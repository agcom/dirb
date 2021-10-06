package jsns

import (
	"errors"
	"github.com/agcom/bs/jsnobjs"
	"github.com/agcom/bs/jsns"
)

type Jsns struct {
	J jsns.Repo
}

func New(j jsns.Repo) *Jsns {
	return &Jsns{j}
}

func (j *Jsns) New(name string, jo jsnobjs.Jsnobj) error {
	err := j.J.New(name, jo)
	return transJsnsError(err)
}

func (j *Jsns) Get(name string) (jsnobjs.Jsnobj, error) {
	jsn, err := j.J.Get(name)
	if err != nil {
		return nil, transJsnsError(err)
	}

	jo, ok := jsn.(jsnobjs.Jsnobj)
	if !ok {
		return nil, ErrJsnButNotJsnobj
	} else {
		return jo, nil
	}
}

func (j *Jsns) Ow(name string, jo jsnobjs.Jsnobj) error {
	err := j.J.Ow(name, jo)
	return transJsnsError(err)
}

func (j *Jsns) Rm(name string) error {
	err := j.J.Rm(name)
	return transJsnsError(err)
}

func (j *Jsns) All() ([]string, error) {
	ns, err := j.J.All()
	return ns, transJsnsError(err)
}

var ErrJsnButNotJsnobj = errors.New("the json is not a json object")

func transJsnsError(err error) error {
	switch err {
	case jsns.ErrExists:
		err = jsnobjs.ErrExists
	case jsns.ErrNotExist:
		err = jsnobjs.ErrNotExist
	}

	return err
}
