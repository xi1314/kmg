package kmgHttp

import (
	"bytes"
	"net/http"
	"testing"

	"net/http/httptest"

	. "github.com/bronze1man/kmg/kmgTest"
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
		MustAddFileToHttpPathToServeMux(mux, "/test/", "test")
		MustAddFileToHttpPathToServeMux(mux, "/test3", "test")
		MustAddFileToHttpPathToServeMux(mux, "test4", "test")
		MustAddFileToHttpPathToServeMux(mux, "/test2/2.html", "test/1.html")
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			panic("should not run to here " + req.URL.String())
		})
		ts := httptest.NewServer(mux)
		defer ts.Close()

		b := MustUrlGetContent(ts.URL + "/test/1.html")
		Equal(b, []byte("1.html"))

		b = MustUrlGetContent(ts.URL + "/test2/2.html")
		Equal(b, []byte("1.html"))

		b = MustUrlGetContent(ts.URL + "/test3/1.html")
		Equal(b, []byte("1.html"))

		b = MustUrlGetContent(ts.URL + "/test4/1.html")
		Equal(b, []byte("1.html"))

		resp, err := http.Get(ts.URL + "/test/2.html")
		Equal(err, nil)
		Equal(resp.StatusCode, 404)
	}

	{
		mux := http.NewServeMux()
		MustAddFileToHttpPathToServeMux(mux, "/test/1.html", "test/1.html")
		ts := httptest.NewServer(mux)
		defer ts.Close()

		b := MustUrlGetContent(ts.URL + "/test/1.html")
		Equal(b, []byte("1.html"))
	}
}
