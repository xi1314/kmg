package kmgErr

import (
	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgLog"
	"runtime/debug"
)

// useful in test
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogErrorWithStack(err error) {
	s := ""
	if err != nil {
		s = err.Error()
	}
	kmgLog.Log("error", s, kmgDebug.GetCurrentStack(1))
}

func LogError(err error) {
	s := ""
	if err != nil {
		s = err.Error()
	}
	kmgLog.Log("error", s)
}

var PrintStack = debug.PrintStack
