package kmgHttp

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestAddParameterToUrl(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	u, err := AddParameterToUrl("http://foo.com/", map[string]string{
		"a": "b",
		"b": "c",
		"c": "d",
	})
	t.Equal(err, nil)
	t.Equal(u, "http://foo.com/?a=b&b=c&c=d")
}

/*
import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestNewUrlByString(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	url, err := NewUrlByString("http://www.google.com")
	t.Equal(nil, err)
	t.Equal("http://www.google.com", url.String())
}
*/
