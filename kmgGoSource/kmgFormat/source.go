package kmgFormat

import (
	"bytes"
	"go/format"
)

func Source(src []byte) ([]byte, error) {
	out, err := format.Source(src)
	if err != nil {
		return out, err
	}
	return bytes.Replace(out, []byte{'\n', '\n'}, []byte{'\n'}, -1), nil
}
