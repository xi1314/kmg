package kmgHttp

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestAlipayContentTypeBug(ot *testing.T) {
	// 正常的(浏览器生成的) Content-Type: application/x-www-form-urlencoded; charset=UTF-8
	// 支付宝的 Content-Type: application/x-www-form-urlencoded; text/html; charset=utf-8
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "http://www.google.com/123", bytes.NewBufferString("a=b&c=d"))
	if err != nil {
		panic(err)
	}
	// 正常的
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	c := NewContextFromHttp(w, req)
	kmgTest.Equal(c.InStr("a"), "b")
	kmgTest.Equal(c.InStr("c"), "d")

	//支付宝的
	req, err = http.NewRequest("POST", "http://www.google.com/123", bytes.NewBufferString("a=b&c=d"))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; text/html; charset=utf-8")
	c = NewContextFromHttp(w, req)
	kmgTest.Equal(c.InStr("a"), "b")
	kmgTest.Equal(c.InStr("c"), "d")
}
