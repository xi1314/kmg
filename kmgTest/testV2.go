package kmgTest

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"reflect"
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
