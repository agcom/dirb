package main

import (
	"github.com/agcom/bs/jsn"
	"path/filepath"
)

var entDirLkp = make(map[string]*jsn.Dir, 3)
var entDirs []*jsn.Dir

func newEnt(name string, aliases []string) {
	entityDir := jsn.NewDir(filepath.Join(rootDir.BinDir().Dir(), name))

	entDirs = append(entDirs, entityDir)

	entDirLkp[name] = entityDir
	for _, alias := range aliases {
		entDirLkp[alias] = entityDir
	}
}

func entDir(name string) *jsn.Dir {
	return entDirLkp[name]
}

const bk = "bk"
const slr = "slr"
const byr = "byr"

func regDefEnts() {
	newEnt("books", []string{"book", bk, "bks"})
	newEnt("sellers", []string{"seller", slr, "slrs"})
	newEnt("buyers", []string{"buyer", byr, "byrs"})
}
