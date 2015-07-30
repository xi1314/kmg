package kmgErr

import (
	"testing"

	. "github.com/bronze1man/kmg/kmgTest"
)

/*
func TestPanic(ot *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		fmt.Println(r)
		debug.PrintStack()
	}()
	f1()
}

func f1() {
	panic(1)
}
*/

func TestPanicToError(ot *testing.T) {
	flag := 1
	err := PanicToError(func() {
		flag = 2
	})
	Equal(flag, 2)
	Equal(err, nil)

	err = PanicToError(func() {
		flag = 3
		panic(nil)
		flag = 4
	})
	Equal(flag, 3)
	Equal(err, nil)

	err = PanicToError(func() {
		panic(1)
		flag = 6
	})
	Equal(flag, 3)
	Ok(err != nil)
	Equal(err.Error(), "1")
}
