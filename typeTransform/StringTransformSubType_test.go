package typeTransform

import "testing"
import (
	"github.com/bronze1man/kmg/kmgTest"
)

type StringTranT1 struct {
	T2 StringTranT2
}
type StringTranT2 string

func TestStringTransformSubType(ot *testing.T) {
	in := &StringTranT1{
		T2: "6",
	}
	err := StringTransformSubType(in, map[string]map[string]string{
		"github.com/bronze1man/kmg/typeTransform.StringTranT2": {
			"6": "Fire",
		},
	})
	kmgTest.Equal(err, nil)
	kmgTest.Equal(in.T2, StringTranT2("Fire"))
}
