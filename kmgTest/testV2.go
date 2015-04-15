package kmgTest

import (
	"fmt"
)

func Ok(expectTrue bool) {
	if !expectTrue {
		panic("ok fail")
	}
}

func Equal(get interface{}, expect interface{}) {
	if isEqual(expect, get) {
		return
	}
	msg := fmt.Sprintf("\tget1: %s\n\texpect2: %s", valueDetail(get), valueDetail(expect))
	panic(msg)
}

type assertPanicType struct{}

//测试此处代码已经panic过了,并且返回panic出来的对象,如果没有panic,此处会panic
//please use github.com/bronze1man/kmgTest
func AssertPanic(f func()) (out interface{}) {
	defer func() {
		out = recover()
		_, ok := out.(assertPanicType)
		if ok {
			panic("should panic")
		}
	}()
	f()
	panic(assertPanicType{})
}

func valueDetail(value interface{}) string {
	stringer, ok := value.(tStringer)
	if ok {
		return fmt.Sprintf("%s (%T) %#v", stringer.String(), value, value)
	} else {
		return fmt.Sprintf("%#v (%T)", value, value)
	}
}

type tStringer interface {
	String() string
}
