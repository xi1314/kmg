package kmgReflect_test

import (
	"reflect"
	"testing"

	"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgTest"
)

type GetAllFieldT1 struct {
	GetAllFieldT3
	*GetAllFieldT4
	B int
}

type GetAllFieldT2 struct {
	A int
	B int
	C int
}

type GetAllFieldT3 struct {
	A int
	B int
	GetAllFieldT2
}
type GetAllFieldT4 struct {
	A int
	D int
}

type GetAllFieldT5 struct {
	GetAllFieldT6
	A int
}
type GetAllFieldT6 int

func TestStructGetAllField(ot *testing.T) {
	t1 := reflect.TypeOf(&GetAllFieldT1{})
	ret := kmgReflect.StructGetAllField(t1)
	kmgTest.Equal(len(ret), 7)
	kmgTest.Equal(ret[0].Name, "GetAllFieldT3")
	kmgTest.Equal(ret[1].Name, "GetAllFieldT4")
	kmgTest.Equal(ret[2].Name, "B")
	kmgTest.Equal(ret[2].Index, []int{2})
	kmgTest.Equal(ret[3].Name, "A")
	kmgTest.Equal(ret[3].Index, []int{0, 0})
	kmgTest.Equal(ret[4].Name, "GetAllFieldT2")
	kmgTest.Equal(ret[5].Name, "C")
	kmgTest.Equal(ret[5].Index, []int{0, 2, 2})
	kmgTest.Equal(ret[6].Name, "D")
	kmgTest.Equal(ret[6].Index, []int{1, 1})

	ret = kmgReflect.StructGetAllField(reflect.TypeOf(&GetAllFieldT5{}))
	kmgTest.Equal(len(ret), 2)

}

func TestStructGetAllFieldMap(ot *testing.T) {
	t1 := reflect.TypeOf(&GetAllFieldT1{})
	ret := kmgReflect.StructGetAllFieldMap(t1)
	kmgTest.Equal(ret["A"].Index, []int{0, 0})
	kmgTest.Equal(ret["B"].Index, []int{2})
	kmgTest.Equal(ret["C"].Index, []int{0, 2, 2})
	kmgTest.Equal(ret["D"].Index, []int{1, 1})
	kmgTest.Equal(len(ret), 7)

	ret = kmgReflect.StructGetAllFieldMap(reflect.TypeOf(&GetAllFieldT5{}))
	kmgTest.Equal(ret["A"].Index, []int{1})
	kmgTest.Equal(len(ret), 2)
}
