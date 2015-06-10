package kmgTextEncoding

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
)

type encodingType string

const (
	Utf8     encodingType = "utf-8"
	ShiftJis encodingType = "shift_jis"
)

var encodingGuessList []encodingType = []encodingType{
	ShiftJis,
	Utf8,
}

//目前只处理了编码是 shift_jis 时的情况
func HttpResponseToUtf8(res *http.Response) (out []byte) {
	body := kmgHttp.MustResponseReadAllBody(res)
	for _, encoding := range encodingGuessList {
		if !isResponseEncodingBy(encoding, res) {
			continue
		}
		if encoding == ShiftJis {
			tReader := transform.NewReader(bytes.NewReader(body), japanese.ShiftJIS.NewDecoder())
			var err error
			out, err = ioutil.ReadAll(tReader)
			if err != nil {
				panic(err)
			}
			return out
		}
		if encoding == Utf8 {
			return body
		}
	}
	panic("[kmgHttp.ResponseToUtf8] Unknown Encoding Type")
	return
}

func isResponseEncodingBy(encoding encodingType, res *http.Response) bool {
	contentType := res.Header.Get("Content-Type")
	charset := getCharsetFromHttpContentType(contentType)
	if charset == string(encoding) {
		return true
	}
	if charset != "" {
		return false
	}
	buf := kmgHttp.MustResponseReadAllBody(res)
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}
	contentType, ok := dom.Find(`meta[http-equiv="content-type"]`).Eq(0).Attr("content")
	if !ok {
		return false
	}
	charset = getCharsetFromHttpContentType(contentType)
	return charset == string(encoding)
}

func getCharsetFromHttpContentType(contentType string) string {
	list := strings.Split(contentType, "charset=")
	if len(list) == 1 {
		return ""
	}
	return strings.ToLower(list[len(list)-1])
}
