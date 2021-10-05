package utils

import (
	"bs/internal/jsns"
	"encoding/json"
	"io"
)

func ReaderToJsn(r io.Reader) (jsns.Jsn, error) {
	dec := json.NewDecoder(r)
	var v jsns.Jsn
	err := dec.Decode(&v)
	return v, err
}

func BytesToJsn(b []byte) (jsns.Jsn, error) {
	var v jsns.Jsn
	err := json.Unmarshal(b, &v)
	return v, err
}

func StrToJsn(s string) (jsns.Jsn, error) {
	return BytesToJsn([]byte(s))
}

func JsnToStr(j jsns.Jsn, indent bool) (string, error) {
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	if indent {
		enc.SetIndent("", "\t")
	}
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		_ = w.CloseWithError(err)
	}()
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
