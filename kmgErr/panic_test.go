package kmgErr

import (
	"fmt"
	"runtime/debug"
	"testing"
)

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
