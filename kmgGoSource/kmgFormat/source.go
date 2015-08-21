package kmgFormat

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
	"go/format"
)

func Source(src []byte) ([]byte, error) {
	src = RemoteEmptyLine(src)
	out, err := format.Source(src)
	if err != nil {
		return out, err
	}
	return out, nil
}

func RemoteEmptyLine(src []byte) []byte {
	r := kmgGoReader.NewReader(src, nil)
	outBuf := &bytes.Buffer{}
	r.ReadAllSpaceWithoutLineBreak()
	if r.IsEof() {
		return []byte{}
	}
	b := r.ReadByte()
	if b != '\n' {
		outBuf.Write(r.BufToCurrent(0))
	}
	for {
		if r.IsEof() {
			return outBuf.Bytes()
		}
		b := r.ReadByte()
		if b == '\n' {
			startPos := r.Pos() - 1 // 包括那个\n
			r.ReadAllSpaceWithoutLineBreak()
			if r.IsEof() {
				outBuf.WriteByte('\n')
				return outBuf.Bytes()
			}
			b := r.ReadByte()
			if b == '\n' {
				r.UnreadByte()
			} else {
				outBuf.Write(r.BufToCurrent(startPos))
			}
		} else {
			outBuf.WriteByte(b)
		}
	}
}
