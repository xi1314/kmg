package kmgIo

import (
	"io"
	"io/ioutil"
)

func MustReadAll(r io.Reader) (b []byte) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}

type StructWriterReaderCloser struct {
	io.Writer
	io.Reader
	io.Closer
}

var NopCloser io.Closer = nopCloser{}

type nopCloser struct{}

func (c nopCloser) Close() (err error) {
	return
}
