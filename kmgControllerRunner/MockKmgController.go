package kmgControllerRunner

import (
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func CallApiWithHttp(uri string, c *kmgHttp.Context) (output string, err error) {
	server := httptest.NewServer(HttpHandler)
	defer server.Close()
	var response *http.Response
	uri = server.URL + uri
	if c.Method == "GET" {
		response, err = http.Get(uri)
	} else {
		postData := url.Values{}
		for key, value := range c.Request {
			postData.Set(key, value)
		}
		response, err = http.PostForm(uri, postData)
	}
	if err != nil {
		return "", err
	}
	_b, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	content := string(_b)
	if err != nil {
		return "", err
	}
	return content, nil
}
