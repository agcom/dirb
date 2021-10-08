package main

import (
	"fmt"
	"github.com/agcom/dirb/bin"
	"github.com/agcom/dirb/jsn"
	"go.uber.org/multierr"
	"path/filepath"
)

type dir jsn.Dir

func newDir(d string) *dir {
	dir := dir(*jsn.NewDir(d))
	return &dir
}

func (d *dir) new(name string, j interface{}) error {
	name = name + ".json"
	return d.jsnDir().New(name, j)
}

func (d *dir) get(name string) (interface{}, error) {
	name = name + ".json"
	return d.jsnDir().Get(name)
}

func (d *dir) over(name string, j interface{}) error {
	name = name + ".json"
	return d.jsnDir().Over(name, j)
}

func (d *dir) rm(name string) error {
	name = name + ".json"
	return d.jsnDir().Rm(name)
}

func (d *dir) all() ([]string, error) {
	ns, err := d.jsnDir().All()
	ons := make([]string, 0, len(ns))
	for _, n := range ns {
		if filepath.Ext(n) == ".json" {
			ons = append(ons, n[:len(n)-len(".json")])
		} else {
			err = multierr.Append(err, fmt.Errorf("missing \".json\" extension in %q", n))
		}
	}

	return ons, err
}

func (d *dir) jsnDir() *jsn.Dir {
	jd := jsn.Dir(*d)
	return &jd
}

func (d *dir) getObj(name string) (map[string]interface{}, error) {
	name = name + ".json"
	return d.jsnDir().GetObj(name)
}

func (d *dir) up(name string, j interface{}) error {
	name = name + ".json"
	return d.jsnDir().Up(name, j)
}

func (d *dir) binDir() *bin.Dir {
	return d.jsnDir().BinDir()
}
