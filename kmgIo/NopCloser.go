package kmgIo

import (
	"io"
)

var NopCloser io.Closer = nopCloser{}

type nopCloser struct{}

func (c nopCloser) Close() (err error) {
	return
}