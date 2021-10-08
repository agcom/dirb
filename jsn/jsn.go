// Package jsn is built on top of bin; focus on json, not binary.
// Meant to be used directly in the CLI application (main function in the root directory of the repository).
package jsn

import (
	"encoding/json"
	"fmt"
	"github.com/agcom/dirb/bin"
	"go.uber.org/multierr"
	"io"
)

func New(path string, j interface{}) error {
	r, w := io.Pipe()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		if err != nil {
			err = fmt.Errorf("failed to encode %v into a json; %w", j, err)
		}
		_ = w.CloseWithError(err)
	}()

	return bin.New(path, r)
}

func Get(path string) (rJ interface{}, rErr error) {
	r, err := bin.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := r.Close()
		if err != nil {
			rErr = multierr.Append(rErr, fmt.Errorf("failed to close binary %q; %w", path, err))
		}
	}()

	j, err := ReaderToJsn(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q into a json; %w", path, err)
	}

	return j, nil
}

func Over(path string, j interface{}) error {
	return bin.Over(path, jsnToReader(j))
}

func Rm(path string) error {
	return bin.Rm(path)
}

func Up(path string, j interface{}) (rErr error) {
	// Early existence check (not vital)
	err := bin.ErrIfNotExist(path)

	lckPath := bin.DefLckPath(path)
	lckFile, err := bin.Lck(lckPath)
	if err != nil {
		return err
	}
	defer func() {
		err := bin.Unlck(lckPath, lckFile)
		if err != nil {
			rErr = multierr.Append(rErr, err)
		}
	}()

	jOld, err := Get(path)
	if err != nil {
		return err
	}
	jNew := mergeJsnRec(jOld, j)

	return bin.OverBare(path, jsnToReader(jNew))
}

func jsnToReader(j interface{}) io.Reader {
	r, w := io.Pipe()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		if err != nil {
			err = fmt.Errorf("failed to encode %v into a json; %w", j, err)
		}
		_ = w.CloseWithError(err)
	}()

	return r
}

func mergeJsnRec(j1, j2 interface{}) interface{} {
	if j1jo, ok := j1.(map[string]interface{}); ok {
		if j2jo, ok := j2.(map[string]interface{}); ok {
			j2 = mergeJsnObjRec(j1jo, j2jo)
		}
	}

	return j2
}

func mergeJsnObjRec(j1, j2 map[string]interface{}) map[string]interface{} {
	r := make(map[string]interface{}, len(j1))
	for k, v1 := range j1 {
		r[k] = v1
	}

	for k, v2 := range j2 {
		if vr, ok := r[k]; ok {
			if vrjo, ok := vr.(map[string]interface{}); ok {
				if v2jo, ok := v2.(map[string]interface{}); ok {
					v2 = mergeJsnObjRec(vrjo, v2jo)
				}
			}
		}

		r[k] = v2
	}

	return r
}
