package kmgHttp

import (
	"bytes"
	"net/http"
	"testing"

	. "github.com/bronze1man/kmg/kmgTest"
	"net/http/httptest"
)

func TestMustRequestToStringCanRead(ot *testing.T) {
	body := bytes.NewReader([]byte("123"))
	req, err := http.NewRequest("POST", "http://foo.com", body)
	Equal(err, nil)
	s1 := MustRequestToStringCanRead(req)
	Equal(s1, "POST / HTTP/1.1\r\n"+
		"Host: foo.com\r\n"+
		"User-Agent: Go 1.1 package http\r\n"+
		"Content-Length: 3\r\n"+
		"\r\n"+
		"123")
	s2 := MustRequestToStringCanRead(req)
	Equal(s1, s2)
	req2 := MustRequestFromString(s1)
	Equal(req2.Host, req.Host)
}

func TestAddFileToHttpPathToServeMux(t *testing.T) {
	{
		mux := http.NewServeMux()
		err := AddFileToHttpPathToServeMux(mux, "/test/", "test")
		Equal(err, nil)
		err = AddFileToHttpPathToServeMux(mux, "/test2", "test")
		Equal(err, nil)
		err = AddFileToHttpPathToServeMux(mux, "/test2/2.html", "test/1.html")
		Equal(err, nil)
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {})
		Equal(err, nil)
		ts := httptest.NewServer(mux)
		defer ts.Close()

		b, err := UrlGetContent(ts.URL + "/test/1.html")
		Equal(err, nil)
		Equal(b, []byte("1.html"))

		b, err = UrlGetContent(ts.URL + "/test2/1.html")
		Equal(err, nil)
		Equal(b, []byte("1.html"))
	}

	{
		mux := http.NewServeMux()
		err := AddFileToHttpPathToServeMux(mux, "/test/1.html", "test/1.html")
		Equal(err, nil)
		ts := httptest.NewServer(mux)
		defer ts.Close()

		b, err := UrlGetContent(ts.URL + "/test/1.html")
		Equal(err, nil)
		Equal(b, []byte("1.html"))
	}
}
