package kmgHttp

import (
	"github.com/bronze1man/kmg/test"
	"testing"
)

func TestNewUrlByString(ot *testing.T) {
	t := test.NewTestTools(ot)
	url, err := NewUrlByString("http://www.google.com")
	t.Equal(nil, err)
	t.Equal("http://www.google.com", url.String())

}
