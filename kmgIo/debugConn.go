package kmgIo

import (
	"fmt"
	"io"
)

//debug a io.ReadWriteCloser
type debugRwc struct {
	io.ReadWriteCloser
	Name string
}

func NewDebugRwc(rwc io.ReadWriteCloser, name string) debugRwc {
	return debugRwc{
		ReadWriteCloser: rwc,
		Name:            name,
	}
}

func (c debugRwc) Write(b []byte) (n int, err error) {
	n, err = c.ReadWriteCloser.Write(b)
	if err != nil {
		fmt.Println("[debugConn]", c.Name, "Write len:", n, "err:", err, "content:", b)
	} else {
		fmt.Println("[debugConn]", c.Name, "Write len:", n, "content:", b)
	}
	return n, err
}

func (c debugRwc) Read(b []byte) (n int, err error) {
	n, err = c.ReadWriteCloser.Read(b)
	if err != nil {
		fmt.Println("[debugConn]", c.Name, "Read inputLen:", len(b), "outputLen:", n, "err:", err, "content:", b[:n])
	} else {
		fmt.Println("[debugConn]", c.Name, "Read inputLen:", len(b), "outputLen:", n, "content:", b[:n])
	}
	return n, err
}

func (c debugRwc) Close() (err error) {
	err = c.ReadWriteCloser.Close()
	fmt.Println("[debugConn]", c.Name, "Close err:", err)
	return err
}
