package kmgControllerRunner

import (
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	. "github.com/bronze1man/kmg/kmgTest"
	"io/ioutil"
	"os"
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
	out := CallApiByHttp(
		"/?n=github.com.bronze1man.kmg.kmgControllerRunner.TestObj.TestFunc&a=10",
		&kmgHttp.Context{Method: "GET"},
	)
	Equal(out, "11")
}

func TestCallApiPost(t *testing.T) {
	RegisterController(TestObj{})
	out := CallApiByHttp(
		"/?n=github.com.bronze1man.kmg.kmgControllerRunner.TestObj.TestFunc",
		&kmgHttp.Context{
			Method: "POST",
			Request: map[string]string{
				"a": "1",
			},
		},
	)
	Equal(out, "2")
}

func TestUploadFile(t *testing.T) {
	testFileRealPath := "/tmp/UFile.md"
	file, err := os.Create(testFileRealPath)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	file.WriteString("hello")
	file.Close()
	RegisterController(TestObj{})
	out := CallApiByHttpWithUploadFile("/?n=github.com.bronze1man.kmg.kmgControllerRunner.TestObj.TestHandleUploadFile",
		&kmgHttp.Context{
			Method: "POST",
			Request: map[string]string{
				"a": "10",
			},
		},
		map[string]string{
			"UFile": "/tmp/UFile.md",
		},
	)
	Equal(out, "UFile.md 10 hello")
}

type TestObj struct{}

func (t TestObj) TestFunc(ctx *kmgHttp.Context) {
	a := ctx.InNum("a")
	ctx.Response = strconv.Itoa(a + i)
}

func (t TestObj) TestHandleUploadFile(ctx *kmgHttp.Context) {
	fileInfo := ctx.InFile("UFile")
	file, err := fileInfo.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	a := ctx.InStr("a")
	ctx.WriteString(fileInfo.Filename)
	ctx.WriteString(" ")
	ctx.WriteString(a)
	ctx.WriteString(" ")
	ctx.WriteString(string(content))
}
