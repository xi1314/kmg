package kmgProfile
import (
	"testing"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgTest"
	"time"
	"net/http"
)

func TestStartProfileOnAddr(ot *testing.T){
	StartProfileOnAddr("abc","127.0.0.1:51001")
	time.Sleep(10*time.Millisecond)
	content:=kmgHttp.MustUrlGetContent("http://127.0.0.1:51001/abc/gc")
	kmgTest.Equal(content,[]byte("SUCCESS"),string(content))

	resp,err:=http.Get("http://127.0.0.1:51001/abc")
	kmgTest.Equal(err,nil)
	kmgTest.Equal(resp.StatusCode,200)

	resp,err=http.Get("http://127.0.0.1:51001/")
	kmgTest.Equal(err,nil)
	kmgTest.Equal(resp.StatusCode,404)

	StartProfileOnAddr("","127.0.0.1:51002")
	time.Sleep(10*time.Millisecond)
	content=kmgHttp.MustUrlGetContent("http://127.0.0.1:51002/gc")
	kmgTest.Equal(content,[]byte("SUCCESS"),string(content))

	resp,err=http.Get("http://127.0.0.1:51002/")
	kmgTest.Equal(err,nil)
	kmgTest.Equal(resp.StatusCode,200)
}
