package internal

import (
	"bs/bins"
	"bs/internal/jsns"
	"go.uber.org/multierr"
	"os"
	"path/filepath"
)

type Repo struct {
	Bks  jsns.Repo
	Slrs jsns.Repo
	Byrs jsns.Repo
}

func NewJsns(bks, slrs, byrs jsns.Repo) *Repo {
	return &Repo{Bks: bks, Slrs: slrs, Byrs: byrs}
}

func NewBins(bks, slrs, byrs bins.Repo) *Repo {
	jbks := jsns.NewJsnsBins(bks)
	jslrs := jsns.NewJsnsBins(slrs)
	jbyrs := jsns.NewJsnsBins(byrs)

	return NewJsns(jbks, jslrs, jbyrs)
}

func NewDirs(bks, slrs, byrs string) *Repo {
	bbks := bins.NewDirBins(bks)
	bslrs := bins.NewDirBins(slrs)
	bbyrs := bins.NewDirBins(byrs)

	return NewBins(bbks, bslrs, bbyrs)
}

func NewDir(d string) (*Repo, error) {
	bks := filepath.Join(d, "books")
	slrs := filepath.Join(d, "sellers")
	byrs := filepath.Join(d, "buyers")

	errbks := os.MkdirAll(bks, 0770)
	errslrs := os.MkdirAll(slrs, 0770)
	errbyrs := os.MkdirAll(byrs, 0770)

	terrs := [3]error{errbks, errslrs, errbyrs}
	var err error = nil
	for _, terr := range terrs {
		if terr != nil {
			err = multierr.Append(err, terr)
		}
	}

	if err != nil {
		return nil, err
	}

	return NewDirs(bks, slrs, byrs), nil
}
