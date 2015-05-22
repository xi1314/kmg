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
	t := kmgTest.NewTestTools(ot)
	name := kmgReflect.GetTypeFullName(reflect.TypeOf(""))
	t.Equal(name, "string")

	name = kmgReflect.GetTypeFullName(reflect.TypeOf(1))
	t.Equal(name, "int")

	name = kmgReflect.GetTypeFullName(reflect.TypeOf(&ta{}))
	t.Equal(name, "github.com/bronze1man/kmg/kmgReflect.ta")

	name = kmgReflect.GetTypeFullName(reflect.TypeOf([]string{}))
	t.Equal(name, "")

}
