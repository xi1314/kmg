package kmgTest

import (
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgControllerRunner"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"net/http"
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

func TestRealCallApi(t *testing.T) {
	kmgControllerRunner.RegisterController(TestObj{})
	http.Handle("/", kmgControllerRunner.HttpHandler)
	err := http.ListenAndServe(":8080", nil)
	kmgConsole.ExitOnErr(err)
}

type TestObj struct{}

func (t TestObj) TestFunc(ctx *kmgHttp.Context) {
	a := ctx.InNum("a")
	ctx.Response = strconv.Itoa(a + i)
}
