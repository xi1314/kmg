package kmgHttp

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestMustRequestToStringCanRead(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	body := bytes.NewReader([]byte("123"))
	req, err := http.NewRequest("POST", "http://foo.com", body)
	t.Equal(err, nil)
	s1 := MustRequestToStringCanRead(req)
	t.Equal(s1, "POST / HTTP/1.1\r\n"+
		"Host: foo.com\r\n"+
		"User-Agent: Go 1.1 package http\r\n"+
		"Content-Length: 3\r\n"+
		"\r\n"+
		"123")
	s2 := MustRequestToStringCanRead(req)
	t.Equal(s1, s2)
	req2 := MustRequestFromString(s1)
	t.Equal(req2.Host, req.Host)
}
