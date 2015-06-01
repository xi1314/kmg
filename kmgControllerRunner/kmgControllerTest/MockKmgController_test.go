package kmgControllerTest

import (
	"github.com/bronze1man/kmg/kmgControllerRunner"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	. "github.com/bronze1man/kmg/kmgTest"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

var i int = 1

func TestMockCallApi(t *testing.T) {
	c := kmgHttp.NewTestContext().
		SetPost().
		SetInStr("a", "1")
	TestObj{}.TestFunc(c)
	Equal(c.GetResponseString(), "2")
}

func TestCallApiGet(t *testing.T) {
	kmgControllerRunner.RegisterController(TestObj{})
	out := CallApiByHttp(
		"/?n=github.com.bronze1man.kmg.kmgControllerRunner.kmgControllerTest.TestObj.TestFunc&a=10",
		kmgHttp.NewTestContext(),
	)
	Equal(out, "11")
}

func TestCallApiPost(t *testing.T) {
	kmgControllerRunner.RegisterController(TestObj{})
	out := CallApiByHttp(
		"/?n=github.com.bronze1man.kmg.kmgControllerRunner.kmgControllerTest.TestObj.TestFunc",
		kmgHttp.NewTestContext().
			SetPost().
			SetInStr("a", "1"))
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
	kmgControllerRunner.RegisterController(TestObj{})
	out := CallApiByHttpWithUploadFile("/?n=github.com.bronze1man.kmg.kmgControllerRunner.kmgControllerTest.TestObj.TestHandleUploadFile",
		kmgHttp.NewTestContext().
			SetPost().
			SetInStr("a", "10"),
		map[string]string{
			"UFile": "/tmp/UFile.md",
		},
	)
	Equal(out, "UFile.md 10 hello")
}

type TestObj struct{}

func (t TestObj) TestFunc(ctx *kmgHttp.Context) {
	a := ctx.InNum("a")
	ctx.WriteString(strconv.Itoa(a + i))
}

func (t TestObj) TestHandleUploadFile(ctx *kmgHttp.Context) {
	fileInfo := ctx.MustInFile("UFile")
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
