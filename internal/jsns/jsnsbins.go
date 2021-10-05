package jsns

import (
	"bs/bins"
	"bs/internal/logs"
	"encoding/json"
	"io"
	"path/filepath"
)

type jsnsBins struct {
	b bins.Repo
}

func (jb *jsnsBins) New(name string, j Jsn) error {
	nb := name + ".json"

	r, w := io.Pipe()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		_ = w.CloseWithError(err)
	}()

	err := jb.b.New(nb, r)
	return transBinsError(err)
}

func (jb *jsnsBins) Get(name string) (Jsn, error) {
	nb := name + ".json"

	r, err := jb.b.Open(nb)
	if err != nil {
		return nil, transBinsError(err)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			logs.Warnf("close bin failed; %v", err)
		}
	}()

	dec := json.NewDecoder(r)
	var j Jsn
	err = dec.Decode(&j)

	if err != nil {
		return nil, err
	} else {
		return j, nil
	}
}

func (jb *jsnsBins) Up(name string, j Jsn) error {
	nb := name + ".json"

	r, w := io.Pipe()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		_ = w.CloseWithError(err)
	}()

	err := jb.b.Ow(nb, r)
	return transBinsError(err)
}

func (jb *jsnsBins) Rm(name string) error {
	nb := name + ".json"
	err := jb.b.Rm(nb)
	return transBinsError(err)
}

func (jb *jsnsBins) All() ([]string, error) {
	ns, err := jb.b.All()
	if err != nil {
		return nil, err
	}

	jns := make([]string, len(ns))[0:]
	for _, n := range ns {
		ext := filepath.Ext(n)
		if ext != "json" {
			logs.Warnf("bin %s is not a json file", n)
		} else {
			jns = append(jns, n)
		}
	}

	return jns, nil
}

func transBinsError(err error) error {
	switch err {
	case bins.ErrExists:
		err = ErrExists
	case bins.ErrBusy:
		err = ErrBusy
	case bins.ErrNotExist:
		err = ErrNotExist
	}

	return err
}
