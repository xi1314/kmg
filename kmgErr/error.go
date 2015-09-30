package kmgErr

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgLog"
)

// useful in test
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogErrorWithStack(err error) {
	if err == nil {
		return
	}
	s := ""
	if err != nil {
		s = err.Error()
	}
	kmgLog.Log("error", s, kmgDebug.GetCurrentStack(1))
}

func LogUserErrorWithStack(err error) {
	if err == nil {
		return
	}
	s := ""
	if err != nil {
		s = err.Error()
	}
	kmgLog.Log("userError", s, kmgDebug.GetCurrentStack(1))
}

func LogError(err error) {
	s := ""
	if err != nil {
		s = err.Error()
	}
	kmgLog.Log("error", s)
}

var PrintStack = debug.PrintStack

type PanicErr struct {
	PanicObj interface{}
}

func (e *PanicErr) Error() string {
	return fmt.Sprintf("%#v", e.PanicObj)
}

// 把panic转换成err返回,
// 没有panic时,返回nil
// panic(nil)会导致返回nil(没有错误)(这个目前没有找到靠谱的方法解决)
func PanicToError(f func()) (err error) {
	defer func() {
		out := recover()
		if out == nil {
			return
		}
		err1, ok := out.(error)
		if ok {
			err = err1
			return
		}
		err = &PanicErr{PanicObj: out}
	}()
	f()
	return nil
}

func PanicToErrorAndLog(f func()) (err error) {
	defer func() {
		out := recover()
		if out == nil {
			err = nil
			return
		}

		err1, ok := out.(error)
		if ok {
			err = err1
			LogErrorWithStack(err)
			return
		}
		err = &PanicErr{PanicObj: out}
		LogErrorWithStack(err)
	}()
	f()
	return nil
}

// 要求这个函数在规定的时间内完成
// 要么在这个时间范围内完成,要么返回错误
// 会在新线程中执行函数
//
func ErrInTime(dur time.Duration, f func()) (err error) {
	return nil
}

func Client(reason string, http_code int) {

}
