package kmgControllerRunner

import (
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	. "github.com/bronze1man/kmg/kmgTest"
	"strconv"
	"testing"
)

var i int = 1

func TestMockCallApi(t *testing.T) {
	c := &kmgHttp.Context{
		Method: "POST",
		Request: map[string]string{
			"a": "1",
		},
	}
	TestObj{}.TestFunc(c)
	Equal(c.Response, "2")
}

func TestCallApiGet(t *testing.T) {
	RegisterController(TestObj{})
	out, err := CallApiWithHttp("/?n=github.com.bronze1man.kmg.kmgControllerRunner.TestObj.TestFunc&a=10", &kmgHttp.Context{Method: "GET"})
	Equal(err, nil)
	Equal(out, "11")
}

func TestCallApiPost(t *testing.T) {
	RegisterController(TestObj{})
	out, err := CallApiWithHttp("/?n=github.com.bronze1man.kmg.kmgControllerRunner.TestObj.TestFunc",
		&kmgHttp.Context{
			Method: "POST",
			Request: map[string]string{
				"a": "1",
			},
		})
	Equal(err, nil)
	Equal(out, "2")
}

type TestObj struct{}

func (t TestObj) TestFunc(ctx *kmgHttp.Context) {
	a := ctx.InNum("a")
	ctx.Response = strconv.Itoa(a + i)
}
