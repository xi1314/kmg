package kmgSession_test

import (
	"github.com/bronze1man/kmg/kmgControllerRunner"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgSession"
	. "github.com/bronze1man/kmg/kmgTest"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func TestNoCookie(t *testing.T) {
	s := kmgSession.GetSessionById("1")
	Equal((s.Id == ""), false)
	s.Set("Name", "Degas")
	s.Set("Age", "25")
	id := s.Id
	s = kmgSession.GetSessionById(id)
	Equal(id, s.Id)
	Equal(s.Get("Name"), "Degas")
	Equal(s.Get("Age"), "25")
}

func TestCookie(t *testing.T) {
	kmgControllerRunner.RegisterController(TestApiObj{})
	uri := "/?n=github.com.bronze1man.kmg.kmgSession_test.TestApiObj.Count"
	server := httptest.NewServer(kmgControllerRunner.HttpHandler)
	defer server.Close()
	var response *http.Response
	uri = server.URL + uri
	cj, _ := cookiejar.New(&cookiejar.Options{})
	client := http.Client{Jar: cj}
	Equal(shareSessionId, "")
	response, err := client.Get(uri)
	Equal(err, nil)
	Equal((shareSessionId == ""), false)
	_sessionId := shareSessionId
	response, _ = client.Get(uri)
	Equal(_sessionId, shareSessionId)
	_b, err := ioutil.ReadAll(response.Body)
	Equal(err, nil)
	response.Body.Close()
	content := string(_b)
	Equal(err, nil)
	Equal(content, shareSessionId)
}

type TestApiObj struct{}

var shareSessionId string

func (t TestApiObj) Count(ctx *kmgHttp.Context) {
	shareSessionId = ctx.Session.Id
	ctx.WriteString(ctx.Session.Id)
}
