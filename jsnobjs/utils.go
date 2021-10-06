package jsnobjs

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
)

func NewJsnobjGenName(r Repo, jo Jsnobj) (string, error) {
	return NewJsnobjGenNameCustom(r, jo, 7, 21, 10000)
}

func NewJsnobjGenNameCustom(r Repo, j Jsnobj, minNameLen, maxNameLen, triesPerLen int) (string, error) {
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
