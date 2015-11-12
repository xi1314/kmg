package kmgReflect_test

import (
	"reflect"
	"testing"

	"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgTest"
)

type ta struct {
}

func TestGetFullName(ot *testing.T) {
	name := kmgReflect.GetTypeFullName(reflect.TypeOf(""))
	kmgTest.Equal(name, "string")

	name = kmgReflect.GetTypeFullName(reflect.TypeOf(1))
	kmgTest.Equal(name, "int")

	name = kmgReflect.GetTypeFullName(reflect.TypeOf(&ta{}))
	kmgTest.Equal(name, "github.com/bronze1man/kmg/kmgReflect_test.ta")

	name = kmgReflect.GetTypeFullName(reflect.TypeOf([]string{}))
	kmgTest.Equal(name, "")

}
