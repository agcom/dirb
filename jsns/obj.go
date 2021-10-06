package jsns

import "errors"

type ObjJsns struct {
	J Repo
}

var ErrJsnButNotJsnObj = errors.New("the json is not a json object")

func (oj *ObjJsns) New(name string, jo JsnObj) error {
	err := oj.J.New(name, jo)
	return transJsnsError(err)
}

func (oj *ObjJsns) Get(name string) (JsnObj, error) {
	jsn, err := oj.J.Get(name)
	if err != nil {
		return nil, transJsnsError(err)
	}

	jo, ok := jsn.(JsnObj)
	if !ok {
		return nil, ErrJsnButNotJsnObj
	} else {
		return jo, nil
	}
}

func (oj *ObjJsns) Ow(name string, jo JsnObj) error {
	err := oj.J.Ow(name, jo)
	return transJsnsError(err)
}

func (oj *ObjJsns) Rm(name string) error {
	err := oj.J.Rm(name)
	return transJsnsError(err)
}

func (oj *ObjJsns) All() ([]string, error) {
	ns, err := oj.J.All()
	return ns, transJsnsError(err)
}

func transJsnsError(err error) error {
	switch err {
	case ErrExists:
		err = ErrExistsObj
	case ErrNotExist:
		err = ErrNotExistObj
	}

	return err
}
