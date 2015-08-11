package kmgGoParser

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

func goSourceRemoveComment(in []byte, filePos *FilePos) (out []byte) {
	r := newReader(in, filePos)
	outBuf := &bytes.Buffer{}
	for {
		switch {
		case r.IsEof():
			return outBuf.Bytes()
		case r.IsMatchAfter([]byte("/*")):
			thisBuf := r.ReadUntilString([]byte("*/"))
			// 保留换行
			commentReader := newReader(thisBuf, nil)
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
			startPos := r.pos
			//fmt.Println("string start "+r.GetFileLineInfo())
			mustReadGoString(r)
			outBuf.Write(r.buf[startPos:r.pos])
		case r.IsMatchAfter([]byte(`'`)):
			startPos := r.pos
			mustReadGoChar(r)
			outBuf.Write(r.buf[startPos:r.pos])
		default:
			outBuf.WriteByte(r.ReadByte())
		}
	}
}

type reader struct {
	buf     []byte //需要读入的数据
	pos     int    //当前位置
	filePos *FilePos
}

func newReader(buf []byte, filePos *FilePos) *reader {
	return &reader{
		buf:     buf,
		filePos: filePos,
	}
}

func (r *reader) IsEof() bool {
	return r.pos >= len(r.buf)
}
func (r *reader) ReadByte() byte {
	out := r.buf[r.pos]
	r.pos++
	return out
}

func (r *reader) NextByte() byte {
	return r.buf[r.pos]
}

func (r *reader) IsMatchAfter(s []byte) bool {
	return bytes.HasPrefix(r.buf[r.pos:], s)
}

// 读取到某个字符,或者读取到结束(该字符会已经被读过)
func (r *reader) ReadUntilByte(b byte) []byte {
	startPos := r.pos
	for {
		if r.IsEof() {
			return r.buf[startPos:]
		}
		if r.ReadByte() == b {
			return r.buf[startPos:r.pos]
		}
	}
}

// 回调返回真的时候,停止读取,(这个回调提到的字符串也包含在内)
func (r *reader) ReadUntilRuneCb(cb func(run rune) bool) []byte {
	startPos := r.pos
	for {
		if r.IsEof() {
			return r.buf[startPos:]
		}
		run, size := utf8.DecodeRune(r.buf[r.pos:])
		r.pos += size
		if cb(run) {
			return r.buf[startPos:r.pos]
		}
	}
}

// 读取到某个字符串,或者读取到结束(该字符串会已经被读过)
func (r *reader) ReadUntilString(s []byte) []byte {
	startPos := r.pos
	for {
		if r.IsEof() {
			return r.buf[startPos:]
		}
		if r.IsMatchAfter(s) {
			r.pos += len(s)
			return r.buf[startPos:r.pos]
		}
		r.pos++
	}
}

func (r *reader) ReadAllSpace() {
	for {
		if r.IsEof() {
			return
		}
		run, size := utf8.DecodeRune(r.buf[r.pos:])
		if !unicode.IsSpace(run) {
			return
		}
		r.pos += size
	}
}

func (r *reader) ReadAllSpaceWithoutLineBreak() {
	for {
		if r.IsEof() {
			return
		}
		run, size := utf8.DecodeRune(r.buf[r.pos:])
		if unicode.IsSpace(run) && run != '\n' {
			r.pos += size
		} else {
			return
		}
	}
}

func (r *reader) ReadRune() rune {
	run, size := utf8.DecodeRune(r.buf[r.pos:])
	r.pos += size
	return run
}

func (r *reader) UnreadRune() rune {
	run, size := utf8.DecodeLastRune(r.buf[:r.pos])
	if size == 0 {
		panic(r.GetFileLineInfo() + " [UnreadRune] last is not valid utf8 code.")
	}
	r.pos -= size
	return run
}

func (r *reader) UnreadByte() {
	r.pos -= 1
}

func (r *reader) MustReadMatch(s []byte) {
	if !r.IsMatchAfter(s) {
		panic(r.GetFileLineInfo() + " [MustReadMatch] not match " + string(s))
	}
	r.pos += len(s)
}

func (r *reader) MustReadWithSize(size int) []byte {
	if r.IsEof() {
		panic(r.GetFileLineInfo() + " unexpect EOF")
	}
	output := r.buf[r.pos : r.pos+size]
	r.pos += size
	return output
}

func (r *reader) GetFileLineInfo() string {
	return r.filePos.PosString(r.pos)
}
