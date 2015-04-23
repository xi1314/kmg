package kmgHttp

import (
	. "github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestAddParameterToUrl(ot *testing.T) {
	u, err := AddParameterToUrl("http://foo.com/", "a", "b")
	Equal(err, nil)
	Equal(u, "http://foo.com/?a=b")

	u, err = AddParameterToUrl("/?n=a&n1=b&n2=c&a=c", "a", "b")
	Equal(err, nil)
	Equal(u, "/?a=c&a=b&n=a&n1=b&n2=c")
}

func TestSetParameterToUrl(ot *testing.T) {
	u, err := SetParameterToUrl("http://foo.com/", "a", "b")
	Equal(err, nil)
	Equal(u, "http://foo.com/?a=b")

	u, err = SetParameterToUrl("/?n=a&n1=b&n2=c&a=c&b=d", "a", "b")
	Equal(err, nil)
	Equal(u, "/?a=b&b=d&n=a&n1=b&n2=c")
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
