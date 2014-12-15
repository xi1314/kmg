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
	fmt.Println("[debugConn]", c.Name, "Write Start len:", len(b))
	n, err = c.ReadWriteCloser.Write(b)
	if err != nil {
		fmt.Println("[debugConn]", c.Name, "Write finish len:", n, "err:", err, "content:", b)
	} else {
		fmt.Println("[debugConn]", c.Name, "Write finish len:", n, "content:", b)
	}
	return n, err
}

func (c debugRwc) Read(b []byte) (n int, err error) {
	fmt.Println("[debugConn]", c.Name, "Read Start len:", len(b))
	n, err = c.ReadWriteCloser.Read(b)
	if err != nil {
		fmt.Println("[debugConn]", c.Name, "Read finish inputLen:", len(b), "outputLen:", n, "err:", err, "content:", b[:n])
	} else {
		fmt.Println("[debugConn]", c.Name, "Read finish inputLen:", len(b), "outputLen:", n, "content:", b[:n])
	}
	return n, err
}

func (c debugRwc) Close() (err error) {
	fmt.Println("[debugConn]", c.Name, "Close start err:", err)
	err = c.ReadWriteCloser.Close()
	fmt.Println("[debugConn]", c.Name, "Close finish err:", err)
	return err
}
