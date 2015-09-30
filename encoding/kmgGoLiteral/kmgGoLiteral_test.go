package kmgGoLiteral

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

type TestType struct {
	A int
	B string
	C []string
}

func TestMarshalToString(ot *testing.T) {
	kmgTest.Equal(MarshalToString(1), "1")
	kmgTest.Equal(MarshalToString("1"), `"1"`)
	kmgTest.Equal(MarshalToString(&TestType{A: 1, B: "1", C: []string{"1", "2"}}),
		`&kmgGoLiteral.TestType{A:1, B:"1", C:[]string{"1", "2"}}`)
}
