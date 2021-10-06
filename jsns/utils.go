package jsns

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"math"
	"strings"
)

func NewJsnGenName(r Repo, j Jsn) (string, error) {
	return NewJsnGenNameCustom(r, j, 7, 21, 10000)
}

func NewJsnGenNameCustom(r Repo, j Jsn, minNameLen, maxNameLen, triesPerLen int) (string, error) {
	if minNameLen > maxNameLen {
		panic("min name len should be less than or equal to max name len")
	}

	name := ""

	for l := minNameLen; l <= maxNameLen; l++ {
		for i := 0; i < triesPerLen; i++ {
			name = genNameLen(l)

			err := r.New(name, j)
			if err != nil && errors.Is(err, ErrExists) {
				continue
			}

			return name, nil
		}
	}

	return "", errors.New(fmt.Sprintf("finding a unique name failed after %v tries", (maxNameLen-minNameLen+1)*triesPerLen))
}

func NewJsnObjGenName(ro RepoObj, jo JsnObj) (string, error) {
	return NewJsnObjGenNameCustom(ro, jo, 7, 21, 10000)
}

func NewJsnObjGenNameCustom(ro RepoObj, jo JsnObj, minNameLen, maxNameLen, triesPerLen int) (string, error) {
	if minNameLen > maxNameLen {
		panic("min name len should be less than or equal to max name len")
	}

	name := ""

	for l := minNameLen; l <= maxNameLen; l++ {
		for i := 0; i < triesPerLen; i++ {
			name = genNameLen(l)

			err := ro.New(name, jo)
			if err != nil && errors.Is(err, ErrExists) {
				continue
			}

			return name, nil
		}
	}

	return "", errors.New(fmt.Sprintf("finding a unique name failed after %v tries", (maxNameLen-minNameLen+1)*triesPerLen))
}

func genNameLen(l int) string {
	if l <= 0 {
		panic("name length shouldn't be negative nor zero")
	}

	rnd := make([]byte, int(math.Ceil(float64(l)*6.0/8.0)))
	_, err := cryptoRand.Read(rnd)
	if err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(rnd)[0:l]
}

type jsnExtMid struct {
	r Repo
}

func JsnExtMid(r Repo) Repo {
	return &jsnExtMid{r}
}

func (jem *jsnExtMid) New(name string, j Jsn) error {
	name = name + ".json"
	return jem.r.New(name, j)
}

func (jem *jsnExtMid) Get(name string) (Jsn, error) {
	name = name + ".json"
	return jem.r.Get(name)
}

func (jem *jsnExtMid) Ow(name string, j Jsn) error {
	name = name + ".json"
	return jem.r.Ow(name, j)
}

func (jem *jsnExtMid) Rm(name string) error {
	name = name + ".json"
	return jem.r.Rm(name)
}

func (jem *jsnExtMid) All() ([]string, error) {
	ns, err := jem.r.All()
	jns := make([]string, len(ns))[:0]
	for _, n := range ns {
		exti := strings.LastIndexByte(n, '.')
		ext := n[exti+1:]
		jn := n[:exti]
		if ext == "json" {
			jns = append(jns, jn)
		} else {
			err = multierr.Append(err, fmt.Errorf("binary \"%s\" is without \".json\" extension in a jsons repository", n))
		}
	}

	return jns, err
}

func WeakUpJsn(r Repo, name string, j Jsn) error {
	jOld, err := r.Get(name)
	if err != nil {
		return err
	}

	jNew := mergeJsnRec(jOld, j)

	return r.Ow(name, jNew)
}

func WeakUpJsnObj(r RepoObj, name string, j JsnObj) error {
	jOld, err := r.Get(name)
	if err != nil {
		return err
	}

	jNew := mergeJsnObjRec(jOld, j)

	return r.Ow(name, jNew)
}

func mergeJsnRec(j1, j2 Jsn) Jsn {
	if j1jo, ok := j1.(JsnObj); ok {
		if j2jo, ok := j2.(JsnObj); ok {
			j2 = mergeJsnObjRec(j1jo, j2jo)
		}
	}

	return j2
}

func mergeJsnObjRec(j1, j2 JsnObj) JsnObj {
	r := make(JsnObj, len(j1))
	for k, v1 := range j1 {
		r[k] = v1
	}

	for k, v2 := range j2 {
		if vr, ok := r[k]; ok {
			if vrjo, ok := vr.(JsnObj); ok {
				if v2jo, ok := v2.(JsnObj); ok {
					v2 = mergeJsnObjRec(vrjo, v2jo)
				}
			}
		}

		r[k] = v2
	}

	return r
}
