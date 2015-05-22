package kmgGoSource

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgTest"
	"golang.org/x/tools/go/types"
	"reflect"
	"testing"
)

type RpcDemo struct {
}

func (s *RpcDemo) PostScoreInt(LbId string, Score int, s1 *RpcDemo) (Info string, err error) {
	return "", nil //ignore the function body first.
}

func TestRpcDemo(ot *testing.T) {
	importPath := "github.com/bronze1man/kmg/kmgGoSource"
	kmgCmd.MustRun("kmg go install " + importPath)
	typ := MustGetGoTypesFromReflect(reflect.TypeOf(&RpcDemo{}))

	typ1, ok := typ.(*types.Pointer)
	kmgTest.Equal(ok, true)

	typ2, ok := typ1.Elem().(*types.Named)
	kmgTest.Equal(ok, true)
	kmgTest.Equal(typ2.NumMethods(), 1)

	obj3 := typ2.Method(0)
	kmgTest.Equal(obj3.Name(), "PostScoreInt")

	typ4, ok := obj3.Type().(*types.Signature)
	kmgTest.Equal(ok, true)
	kmgTest.Equal(typ4.Params().Len(), 3)
	kmgTest.Equal(typ4.Results().Len(), 2)

	for _, testCase := range []struct {
		Type   types.Type
		Expect string
	}{
		{typ4.Params().At(0).Type(), "string"},
		{typ4.Params().At(1).Type(), "int"},
		{typ4.Params().At(2).Type(), "*RpcDemo"},
		{typ4.Results().At(0).Type(), "string"},
		{typ4.Results().At(1).Type(), "error"},
	} {
		typS, importPathList := MustWriteGoTypes("github.com/bronze1man/kmg/kmgGoSource", testCase.Type)
		kmgTest.Equal(typS, testCase.Expect)
		kmgTest.Equal(len(importPathList), 0)
	}

	typS, importPathList := MustWriteGoTypes("github.com/bronze1man/kmg/kmgTest", typ4.Params().At(2).Type())
	kmgTest.Equal(typS, "*kmgGoSource.RpcDemo")
	kmgTest.Equal(importPathList, []string{"github.com/bronze1man/kmg/kmgGoSource"})

	methodSet := types.NewMethodSet(typ)
	kmgTest.Equal(methodSet.Len(), 1)
	kmgTest.Equal(methodSet.At(0).Obj().Name(), "PostScoreInt")

}
