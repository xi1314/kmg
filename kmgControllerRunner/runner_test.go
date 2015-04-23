package kmgControllerRunner

import (
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	. "github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestRegisterController(t *testing.T) {
	clearController()
	obj := TestRegisterControllerT{}
	RegisterController(obj)
	Equal(len(controllerObjMap), 2)
	ctx := &kmgHttp.Context{}

	testCallMethod = []string{}
	controllerObjMap["github.com.bronze1man.kmg.kmgControllerRunner.TestRegisterControllerT.WorkPage1"](ctx)
	Equal(testCallMethod, []string{"WorkPage1"})

	controllerObjMap["github.com.bronze1man.kmg.kmgControllerRunner.TestRegisterControllerT.WorkPage2"](ctx)
	Equal(testCallMethod, []string{"WorkPage1", "WorkPage2"})
}

var testCallMethod = []string{}

type TestRegisterControllerT struct{}

func (t TestRegisterControllerT) NotWorkApi() {
	testCallMethod = append(testCallMethod, "NotWorkApi")
}
func (t TestRegisterControllerT) WorkPage1(ctx *kmgHttp.Context) {
	testCallMethod = append(testCallMethod, "WorkPage1")
}
func (t TestRegisterControllerT) WorkPage2(ctx *kmgHttp.Context) {
	testCallMethod = append(testCallMethod, "WorkPage2")
}

func (t TestRegisterControllerT) workPage2(ctx *kmgHttp.Context) {
	testCallMethod = append(testCallMethod, "workPage2")
}
