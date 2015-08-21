package kmgIo

import (
	"github.com/bronze1man/kmg/errors"
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

func MultiErrorHandle(fs ...func() error) error {
	var errS string
	for _, f := range fs {
		err1 := f()
		if err1 != nil {
			errS += "[" + err1.Error() + "] "
		}
	}
	return errors.New(errS)
}
