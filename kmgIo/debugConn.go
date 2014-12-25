package kmgIo

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
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

type sumSizeRwc struct {
	io.ReadWriteCloser
	Name       string
	startTime  time.Time
	writeBytes uint64
	readBytes  uint64
	readNum    uint64
	writeNum   uint64
	hasClose   bool
}

func (c *sumSizeRwc) Write(b []byte) (n int, err error) {
	n, err = c.ReadWriteCloser.Write(b)
	if n > 0 {
		atomic.AddUint64(&c.writeBytes, uint64(n))
		atomic.AddUint64(&c.writeNum, 1)
	}
	return n, err
}

func (c *sumSizeRwc) Read(b []byte) (n int, err error) {
	n, err = c.ReadWriteCloser.Read(b)
	if n > 0 {
		atomic.AddUint64(&c.readBytes, uint64(n))
		atomic.AddUint64(&c.readNum, 1)
	}
	return n, err
}

func (c *sumSizeRwc) Close() (err error) {
	err = c.ReadWriteCloser.Close()
	if !c.hasClose {
		fmt.Printf("[sumSizeRwc] [%s] read[bytes:%s num:%d] write[bytes:%s num:%d] duration:%s\n",
			c.Name, FmtByteNum(int(c.readBytes)), c.readNum, FmtByteNum(int(c.writeBytes)), c.writeNum,
			time.Since(c.startTime))
		c.hasClose = true
	}
	return err
}

func NewSumSizeRwc(rwc io.ReadWriteCloser, name string) io.ReadWriteCloser {
	return &sumSizeRwc{
		ReadWriteCloser: rwc,
		Name:            name,
		startTime:       time.Now(),
	}
}
