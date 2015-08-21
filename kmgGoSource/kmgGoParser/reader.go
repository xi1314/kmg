package kmgGoParser

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
)

func goSourceRemoveComment(in []byte, filePos *kmgGoReader.FilePos) (out []byte) {
	r := kmgGoReader.NewReader(in, filePos)
	outBuf := &bytes.Buffer{}
	for {
		switch {
		case r.IsEof():
			return outBuf.Bytes()
		case r.IsMatchAfter([]byte("/*")):
			thisBuf := r.ReadUntilString([]byte("*/"))
			// 保留换行
			commentReader := kmgGoReader.NewReader(thisBuf, nil)
			for {
				if commentReader.IsEof() {
					break
				}
				b := commentReader.ReadByte()
				if b == '\n' {
					outBuf.WriteByte('\n')
				} else {
					outBuf.WriteByte(' ')
				}
			}
		case r.IsMatchAfter([]byte("//")):
			thisBuf := r.ReadUntilByte('\n')
			outBuf.Write(bytes.Repeat([]byte{' '}, len(thisBuf)-1))
			outBuf.WriteByte('\n')
		case r.IsMatchAfter([]byte(`"`)) || r.IsMatchAfter([]byte("`")):
			startPos := r.Pos()
			//fmt.Println("string start "+r.GetFileLineInfo())
			mustReadGoString(r)
			outBuf.Write(r.BufToCurrent(startPos))
		case r.IsMatchAfter([]byte(`'`)):
			startPos := r.Pos()
			mustReadGoChar(r)
			outBuf.Write(r.BufToCurrent(startPos))
		default:
			outBuf.WriteByte(r.ReadByte())
		}
	}
}
