package kmgBinary

import (
	"fmt"
	"io"
)

// 只能写入长度小于255的字符串.
func MustWriteString255(w io.Writer, s string) {
	if len(s) > 255 {
		panic(fmt.Errorf("[MustWriteString255] len(s)[%d]>255", len(s)))
	}
	_, err := w.Write([]byte{byte(len(s))})
	if err != nil {
		panic(err)
	}
	_, err = w.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

// 读取用 MustWriteString255 写入的东西
func ReadString255(r io.Reader) (s string, err error) {
	buf := make([]byte, 256)
	_, err = io.ReadFull(r, buf[:1])
	if err != nil {
		return "", err
	}
	len := int(buf[0])
	_, err = io.ReadFull(r, buf[:len])
	if err != nil {
		return "", err
	}
	return string(buf[:len]), nil
}
