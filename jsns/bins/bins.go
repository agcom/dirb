package bins

import (
	"encoding/json"
	"fmt"
	"github.com/agcom/bs/bins"
	"github.com/agcom/bs/jsns"
	"go.uber.org/multierr"
	"io"
)

type Bins struct {
	B bins.Repo
}

func New(b bins.Repo) *Bins {
	return &Bins{b}
}

func (b *Bins) New(name string, j jsns.Jsn) error {
	r, w := io.Pipe()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		_ = w.CloseWithError(err)
	}()

	err := b.B.New(name, r)
	return transBinsError(err)
}

func (b *Bins) Get(name string) (rj jsns.Jsn, rErr error) {
	r, err := b.B.Open(name)
	if err != nil {
		return nil, transBinsError(err)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			rErr = multierr.Append(rErr, NewErrBinClose(name, err))
		}
	}()

	dec := json.NewDecoder(r)
	var j jsns.Jsn
	err = dec.Decode(&j)

	if err != nil {
		return nil, err
	} else {
		return j, nil
	}
}

func (b *Bins) Ow(name string, j jsns.Jsn) error {
	r, w := io.Pipe()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		_ = w.CloseWithError(err)
	}()

	err := b.B.Ow(name, r)
	return transBinsError(err)
}

func (b *Bins) Rm(name string) error {
	err := b.B.Rm(name)
	return transBinsError(err)
}

func (b *Bins) All() ([]string, error) {
	ns, err := b.B.All()
	return ns, transBinsError(err)
}

type ErrBinClose struct {
	name  string
	cause error
}

func NewErrBinClose(name string, cause error) *ErrBinClose {
	return &ErrBinClose{
		name:  name,
		cause: cause,
	}
}

func (e *ErrBinClose) Error() string {
	return fmt.Sprintf("failed to close bin \"%s\"; %v", e.name, e.cause)
}

func (e *ErrBinClose) Unwrap() error {
	return e.cause
}

func transBinsError(err error) error {
	switch err {
	case bins.ErrExists:
		err = jsns.ErrExists
	case bins.ErrNotExist:
		err = jsns.ErrNotExist
	}

	return err
}
