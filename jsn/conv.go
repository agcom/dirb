package jsn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func ReaderToJsn(r io.Reader) (interface{}, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	var j interface{}
	err := dec.Decode(&j)
	if err != nil {
		return nil, err
	}

	// Check whether the stream ended; ensure that the whole stream contained exactly a single json.
	_, err = dec.Token()
	if err != nil && errors.Is(err, io.EOF) {
		// Stream ended
		return j, nil
	} else {
		// Stream not ended; can or can't proceed decoding.
		return nil, fmt.Errorf("json already ended, but got more input")
	}
}

func ReaderToJsnObj(r io.Reader) (map[string]interface{}, error) {
	j, err := ReaderToJsn(r)
	if err != nil {
		return nil, err
	}

	if jo, ok := j.(map[string]interface{}); ok {
		return jo, nil
	} else {
		return nil, fmt.Errorf("\"%v\" is a json, but not a json object", j)
	}
}

func ByteSliceToJsn(b []byte) (interface{}, error) {
	return ReaderToJsn(bytes.NewReader(b))
}

func ByteSliceToJsnObj(b []byte) (map[string]interface{}, error) {
	return ReaderToJsnObj(bytes.NewReader(b))
}

func StrToJsn(s string) (interface{}, error) {
	return ByteSliceToJsn([]byte(s))
}

func StrToJsnObj(s string) (map[string]interface{}, error) {
	return ByteSliceToJsnObj([]byte(s))
}

func JsnToReader(j interface{}) io.Reader {
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		_ = w.CloseWithError(err)
	}()

	return r
}

func JsnObjToReader(jo map[string]interface{}) io.Reader {
	return jsnToReader(jo)
}

func JsnToByteSlice(j interface{}) ([]byte, error) {
	bs, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func JsnObjToByteSlice(jo map[string]interface{}) ([]byte, error) {
	return JsnToByteSlice(jo)
}

func JsnToStr(j interface{}) (string, error) {
	bs, err := JsnToByteSlice(j)
	return string(bs), err
}

func JsnObjToStr(jo map[string]interface{}) (string, error) {
	return JsnToStr(jo)
}
