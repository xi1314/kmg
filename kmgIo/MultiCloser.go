package kmgIo

import (
	"io"
)

func MultiCloser(closers ...io.Closer) io.Closer {
	c := make([]io.Closer, len(closers))
	copy(c, closers)
	return multiCloser(c)
}

type multiCloser []io.Closer

func (c multiCloser) Close() (err error) {
	for _, closer := range c {
		err1 := closer.Close()
		if err1 != nil {
			err = err1
		}
	}
	return err
}
