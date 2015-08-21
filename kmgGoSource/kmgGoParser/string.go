package kmgGoParser

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgGoSource/kmgGoReader"
	"strconv"
)

// 可以读这两种 "abc" `abc`
func mustReadGoString(r *kmgGoReader.Reader) (output []byte) {
	b := r.ReadByte()
	if b == '"' {
		buf := &bytes.Buffer{}
		for {
			if r.IsEof() {
				panic(r.GetFileLineInfo() + " unexcept EOF " + buf.String())
			}
			b := r.ReadByte()
			switch b {
			case '"':
				return buf.Bytes()
			case '\\':
				handleSlashInGoChar(r, buf)
			default:
				buf.WriteByte(b)
			}
		}
	} else if b == '`' {
		return r.ReadUntilByte('`')
	} else {
		panic(fmt.Errorf("%s unexcept byte %d '%s'", r.GetFileLineInfo(), b, string(rune(b))))
	}
}

// 可以读 '"' 这种,返回这个东西实际占的字节数据
func mustReadGoChar(r *kmgGoReader.Reader) []byte {
	b := r.ReadByte()
	if b != '\'' {
		panic(r.GetFileLineInfo() + " unexcept byte " + strconv.Itoa(int(b)))
	}
	buf := &bytes.Buffer{}
	run := r.ReadRune()
	if run == '\\' {
		handleSlashInGoChar(r, buf)
	} else {
		buf.WriteRune(run)
	}
	b = r.ReadByte()
	if b != '\'' {
		panic(r.GetFileLineInfo() + " unexcept byte " + strconv.Itoa(int(b)))
	}
	return buf.Bytes()
}

// 此处已经读过了 /
func handleSlashInGoChar(r *kmgGoReader.Reader, buf *bytes.Buffer) {
	if r.IsEof() {
		panic(r.GetFileLineInfo() + " unexcept EOF")
	}
	b := r.ReadByte()
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7':
		r.UnreadByte()
		octal := r.MustReadWithSize(3)
		b64, err := strconv.ParseUint(string(octal), 8, 8)
		if err != nil {
			panic(r.GetFileLineInfo() + " " + err.Error())
		}
		buf.WriteByte(byte(b64))
	case 'x':
		octal := r.MustReadWithSize(2)
		b64, err := strconv.ParseUint(string(octal), 16, 8)
		if err != nil {
			panic(r.GetFileLineInfo() + " " + err.Error())
		}
		buf.WriteByte(byte(b64))
	case 'u':
		octal := r.MustReadWithSize(4)
		b64, err := strconv.ParseUint(string(octal), 16, 16)
		if err != nil {
			panic(r.GetFileLineInfo() + " " + err.Error())
		}
		buf.WriteRune(rune(b64))
	case 'U':
		octal := r.MustReadWithSize(8)
		b64, err := strconv.ParseUint(string(octal), 16, 32)
		if err != nil {
			panic(r.GetFileLineInfo() + " " + err.Error())
		}
		buf.WriteRune(rune(b64))
	case 'a':
		buf.WriteByte('\a')
	case 'b':
		buf.WriteByte('\b')
	case 'f':
		buf.WriteByte('\f')
	case 'n':
		buf.WriteByte('\n')
	case 'r':
		buf.WriteByte('\r')
	case 't':
		buf.WriteByte('\t')
	case 'v':
		buf.WriteByte('\v')
	case '\\':
		buf.WriteByte('\\')
	case '\'':
		buf.WriteByte('\'')
	case '"':
		buf.WriteByte('"')
	default:
		panic(r.GetFileLineInfo() + " unexcept byte " + strconv.Itoa(int(b)))
	}
}
