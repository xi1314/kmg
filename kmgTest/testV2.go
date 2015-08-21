package kmgTest

import (
	"fmt"
	"reflect"

	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgReflect"
)

func Ok(expectTrue bool, objList ...interface{}) {
	if !expectTrue {
		if len(objList) == 0 {
			panic("ok fail")
		} else {
			panic("ok fail\n" + kmgDebug.Sprintln(objList...))
		}
	}
}

func Equal(get interface{}, expect interface{}, objList ...interface{}) {
	if isEqual(expect, get) {
		return
	}
	var msg string
	if len(objList) == 0 {
		msg = fmt.Sprintf("\tget1: %s\n\texpect2: %s", valueDetail(get), valueDetail(expect))
	} else {
		msg = fmt.Sprintf("\tget1: %s\n\texpect2: %s\n%s", valueDetail(get), valueDetail(expect), kmgDebug.Sprintln(objList...))
	}
	panic(msg)
}

type assertPanicType struct{}

//测试此处代码已经panic过了,并且返回panic出来的对象,如果没有panic,此处会panic
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

func isEqual(a interface{}, b interface{}) bool {
	//不要加a==b这种快速通道,会出现 panic comparing uncomparable type []models.RoleId 的问题
	if reflect.DeepEqual(a, b) {
		return true
	}
	rva := reflect.ValueOf(a)
	rvb := reflect.ValueOf(b)
	//every nil is same stuff...
	if kmgReflect.IsNil(rva) && kmgReflect.IsNil(rvb) {
		return true
	}
	return false
}
