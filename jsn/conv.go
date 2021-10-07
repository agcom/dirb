package jsn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func ReaderToJsn(r io.Reader) (interface{}, error) {
	dec := json.NewDecoder(r)
	var j interface{}
	err := dec.Decode(&j)
	if err != nil {
		return nil, fmt.Errorf("failed to decode into a json; %w", err)
	}

	// Check whether the stream ended; ensure that the whole stream contained exactly a single json.
	_, err = dec.Token()
	if err != nil && errors.Is(err, io.EOF) {
		// Stream ended
		return j, nil
	} else {
		// Stream not ended; can or can't proceed decoding.
		return nil, fmt.Errorf("failed to decode into a json; json already ended, but got more input")
	}
}
