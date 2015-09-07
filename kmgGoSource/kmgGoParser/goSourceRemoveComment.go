package kmgGoParser

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
)

type noGrowBuf struct {
	buf []byte
	pos int
}

func (buf *noGrowBuf) WriteByte(b byte) {
	buf.buf[buf.pos] = b
	buf.pos++
}
func (buf *noGrowBuf) Write(b []byte) {
	copy(buf.buf[buf.pos:], b)
	buf.pos += len(b)
}
func goSourceRemoveComment(in []byte, filePos *kmgGoReader.FilePos) (out []byte) {
	r := kmgGoReader.NewReader(in, filePos)
	buf := &noGrowBuf{
		buf: make([]byte, len(in)),
		pos: 0,
	}
	for {
		if r.IsEof() {
			return buf.buf
		}
		b := r.ReadByte()
		switch b {
		case '/':
			r.UnreadByte()
			if r.IsMatchAfter(tokenSlashStar) {
				thisBuf := r.ReadUntilString(tokenStarSlash)
				// 保留换行
				commentReader := kmgGoReader.NewReader(thisBuf, nil)
				for {
					if commentReader.IsEof() {
						break
					}
					b := commentReader.ReadByte()
					if b == '\n' {
						buf.WriteByte('\n')
					} else {
						buf.WriteByte(' ')
					}
				}
			} else if r.IsMatchAfter(tokenDoubleSlash) {
				thisBuf := r.ReadUntilByte('\n')
				buf.Write(bytes.Repeat([]byte{' '}, len(thisBuf)-1))
				buf.WriteByte('\n')
			} else {
				buf.WriteByte(r.ReadByte())
			}
		case '"', '`':
			r.UnreadByte()
			startPos := r.Pos()
			//fmt.Println("string start "+r.GetFileLineInfo())
			mustReadGoString(r)
			buf.Write(r.BufToCurrent(startPos))
		case '\'':
			r.UnreadByte()
			startPos := r.Pos()
			mustReadGoChar(r)
			buf.Write(r.BufToCurrent(startPos))
		default:
			buf.WriteByte(b)
		}
		/*
			switch {
			case r.IsEof():
				return outBuf.Bytes()
			case r.IsMatchAfter(tokenSlashStar):
				thisBuf := r.ReadUntilString(tokenStarSlash)
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
			case r.IsMatchAfter(tokenDoubleSlash):
				thisBuf := r.ReadUntilByte('\n')
				outBuf.Write(bytes.Repeat([]byte{' '}, len(thisBuf)-1))
				outBuf.WriteByte('\n')
			case r.IsMatchAfter(tokenDoubleQuate) || r.IsMatchAfter(tokenGraveAccent):
				startPos := r.Pos()
				//fmt.Println("string start "+r.GetFileLineInfo())
				mustReadGoString(r)
				outBuf.Write(r.BufToCurrent(startPos))
			case r.IsMatchAfter(tokenSingleQuate):
				startPos := r.Pos()
				mustReadGoChar(r)
				outBuf.Write(r.BufToCurrent(startPos))
			default:
				outBuf.WriteByte(r.ReadByte())
			}
		*/
	}
}
