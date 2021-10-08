package jsn

import (
	"github.com/agcom/dirb/bin"
	"path/filepath"
)

type Dir bin.Dir

func NewDir(d string) *Dir {
	dir := Dir(*bin.NewDir(d))
	return &dir
}

func (d *Dir) New(name string, j interface{}) error {
	path := d.Path(name)
	return New(path, j)
}

func (d *Dir) Get(name string) (interface{}, error) {
	path := d.Path(name)
	return Get(path)
}

func (d *Dir) Over(name string, j interface{}) error {
	path := d.Path(name)
	return Over(path, j)
}

func (d *Dir) Rm(name string) error {
	path := d.Path(name)
	return Rm(path)
}

func (d *Dir) All() ([]string, error) {
	return d.BinDir().All()
}

func (d *Dir) BinDir() *bin.Dir {
	bd := bin.Dir(*d)
	return &bd
}

func (d *Dir) GetObj(name string) (map[string]interface{}, error) {
	path := d.Path(name)
	return GetObj(path)
}

func (d *Dir) Up(name string, j interface{}) error {
	path := d.Path(name)
	return Up(path, j)
}

func (d *Dir) dir() string {
	return string(*d.BinDir())
}

func (d *Dir) Path(name string) string {
	return filepath.Join(d.dir(), name)
}
