package kmgIo

import (
	"io"
)

//把一个rwc转换成另一个rwc
type RwcWrapper interface {
	RwcWrap(in io.ReadWriteCloser)(out io.ReadWriteCloser,err error)
}

type RwcWrapperFunc func(in io.ReadWriteCloser)(out io.ReadWriteCloser,err error)

func (f RwcWrapperFunc)RwcWrap(in io.ReadWriteCloser)(out io.ReadWriteCloser,err error){
	return f(in)
}