package kmgTest

import (
	"fmt"
	"reflect"
	"testing"

	"strings"

	"github.com/bronze1man/kmg/kmgDebug" //TODO 移除这个依赖?
)

// @deprecated
type TestingTB interface {
	FailNow()
}

// @deprecated
type TestingTBAware interface {
	SetTestingTB(T TestingTB)
}

// @deprecated
type TestTools struct {
	TestingTB
}

// @deprecated
func NewTestTools(T TestingTB) *TestTools {
	return &TestTools{TestingTB: T}
}

// @deprecated
func (tools *TestTools) SetTestingTB(T TestingTB) {
	tools.TestingTB = T
}

// @deprecated
func (tools *TestTools) Ok(expectTrue bool) {
	if !expectTrue {
		tools.assertFail("ok fail", 2)
	}
	return
}

// @deprecated
func (tools *TestTools) Equal(get interface{}, expect interface{}) {
	if isEqual(expect, get) {
		return
	}
	msg := fmt.Sprintf("\texpect2: %#v (%s) (%T)\n\tget1: %#v (%s) (%T)", expect, expect, expect, get, get, get)
	if eGet, ok := get.(error); ok {
		msg += "\ngetError: " + eGet.Error()
	}
	tools.assertFail(msg, 2)
}

// @deprecated
func (tools *TestTools) EqualMsg(get interface{}, expect interface{}, format string, args ...interface{}) {
	if isEqual(expect, get) {
		return
	}
	tools.assertFail(fmt.Sprintf(`%s
get1:%#v (%T)
expect2:%#v (%T)`, fmt.Sprintf(format, args...), get, get, expect, expect), 2)
}
func (tools *TestTools) GetTestingT() *testing.T {
	return tools.TestingTB.(*testing.T)
}

//an easy way to Printf(avoid of import fmt,and remove it ...)should for Debug only
func (tools *TestTools) Printf(format string, objs ...interface{}) (n int, err error) {
	return fmt.Printf(format, objs...)
}

//an easy way to Println (avoid of import fmt,and remove it ...)should for Debug only
func (tools *TestTools) Println(objs ...interface{}) (n int, err error) {
	return fmt.Println(objs...)
}
func (tools *TestTools) Fatalf(format string, objs ...interface{}) {
	tools.Printf("\n%s\n\n%s\n", fmt.Sprintf(format, objs...),
		kmgDebug.GetCurrentStack(1).ToString())
	tools.TestingTB.FailNow()
}
func (tools *TestTools) assertFail(msg string, skip int) {
	tools.Printf(`----------------------------------
%s

%s----------------------------------
`, msg, kmgDebug.GetCurrentStack(skip).ToString())
	tools.TestingTB.FailNow()
}

// @deprecated
//use T for error report,automatic call every function start with Test
//function start with Test must not have input arguments,output arguments.
func TestWarpper(T TestingTB, testObject TestingTBAware) {
	testObject.SetTestingTB(T)
	tov := reflect.ValueOf(testObject)
	//tov := reflect.Indirect(reflect.ValueOf(testObject))
	tot := tov.Type()
	for i := 0; i < tov.NumMethod(); i++ {
		tm := tot.Method(i)
		if !strings.HasPrefix(tm.Name, "Test") {
			continue
		}
		tmt := tm.Type
		//no argument
		if tmt.NumIn() != 1 {
			fmt.Printf("[kmgTest.TestWarpper] Testfunction:%s should not have input argument\n", tm.Name)
			T.FailNow()
			return
		}
		if tmt.NumOut() != 0 {
			fmt.Printf("[kmgTest.TestWarpper] Testfunction:%s should not have output argument\n", tm.Name)
			T.FailNow()
			return
		}
		tov.Method(i).Call([]reflect.Value{})
	}
}

// @deprecated
func EqualMsg(get interface{}, expect interface{}, objList ...interface{}) {
	if isEqual(expect, get) {
		return
	}
	msg := fmt.Sprintf("\tget1: %s\n\texpect2: %s\n%s", valueDetail(get), valueDetail(expect), kmgDebug.Sprintln(objList...))
	panic(msg)
}
