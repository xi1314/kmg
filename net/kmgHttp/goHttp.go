package kmgHttp

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func NewRequestFromByte(r []byte) (req *http.Request, err error) {
	return http.ReadRequest(bufio.NewReader(bytes.NewReader(r)))
}

//sometimes it is hard to remember how to get response from bytes ...
func NewResponseFromBytes(r []byte) (resp *http.Response, err error) {
	return http.ReadResponse(bufio.NewReader(bytes.NewBuffer(r)), &http.Request{})
}

//sometimes it is hard to remember how to dump response to bytes
func DumpResponseToBytes(resp *http.Response) (b []byte, err error) {
	return httputil.DumpResponse(resp, true)
}

func ResponseReadAllBody(resp *http.Response) (b []byte, err error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func UrlGetContent(url string) (b []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	return ResponseReadAllBody(resp)
}
