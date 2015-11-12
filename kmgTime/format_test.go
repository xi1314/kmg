package kmgTime

import (
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestFormat(ot *testing.T) {
	t, err := time.Parse(AppleJsonFormat, "2014-04-16 18:26:18 Etc/GMT")
	kmgTest.Equal(err, nil)
	kmgTest.Ok(t.Equal(MustFromMysqlFormat("2014-04-16 18:26:18")))
}
