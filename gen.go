package main

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/agcom/bs/bin"
	"github.com/agcom/bs/jsn"
	"math"
)

func newJsnGenName(d *jsn.Dir, j interface{}) (string, error) {
	return newJsnGenNameCustom(d, j, 7, 21, 10000)
}

func newJsnGenNameCustom(d *jsn.Dir, j interface{}, minNameLen, maxNameLen, triesPerLen int) (string, error) {
	if minNameLen > maxNameLen {
		panic(fmt.Sprintf("the minimum name length %d is more than the maximum name length %d", minNameLen, maxNameLen))
	} else if triesPerLen <= 0 {
		panic(fmt.Sprintf("non-positive tries per length %d", triesPerLen))
	}

	name := ""
	for l := minNameLen; l <= maxNameLen; l++ {
		for i := 0; i < triesPerLen; i++ {
			name = genNameLen(l)

			err := d.New(name, j)
			if err != nil {
				if _, ok := err.(*bin.ErrExists); ok {
					continue
				} else {
					return "", err
				}
			}

			return name, nil
		}
	}

	return "", fmt.Errorf("failed to find a unique name after %v tries", (maxNameLen-minNameLen+1)*triesPerLen)
}

func genNameLen(l int) string {
	if l <= 0 {
		panic(fmt.Sprintf("the name length %d is not positive", l))
	}

	rnd := make([]byte, int(math.Ceil(float64(l)*6.0/8.0)))
	_, err := cryptoRand.Read(rnd)
	if err != nil {
		panic(fmt.Errorf("cryptocurrency random number generator failed; %w", err))
	}

	return base64.RawURLEncoding.EncodeToString(rnd)[0:l]
}
