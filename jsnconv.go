package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func strToJsnObj(s string) (map[string]interface{}, error) {
	return bytesToJsnObj([]byte(s))
}

func strToJsn(s string) (interface{}, error) {
	return bytesToJsn([]byte(s))
}

func bytesToJsn(b []byte) (interface{}, error) {
	return readerToJsn(bytes.NewReader(b))
}

func readerToJsn(r io.Reader) (interface{}, error) {
	dec := json.NewDecoder(r)
	var j interface{}
	err := dec.Decode(&j)
	if err != nil {
		return nil, fmt.Errorf("failed to decode; %w", err)
	}

	return j, nil
}

func readerToJsnObj(r io.Reader) (map[string]interface{}, error) {
	j, err := readerToJsn(r)
	if err != nil {
		return nil, err
	}

	if jo, ok := j.(map[string]interface{}); ok {
		return jo, nil
	} else {
		return nil, fmt.Errorf("\"%v\" is a json, but not a json object", j)
	}
}

func bytesToJsnObj(b []byte) (map[string]interface{}, error) {
	return readerToJsnObj(bytes.NewReader(b))
}

func jsnObjToStr(jo map[string]interface{}, tabIndent bool) (string, error) {
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	if tabIndent {
		enc.SetIndent("", "\t")
	}
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(jo)
		_ = w.CloseWithError(err)
	}()
	b, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to encode into a json string; %w", err)
	}

	return string(b), nil
}
