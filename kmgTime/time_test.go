package kmgTime

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestMustFromMysqlFormatDefaultTZ(ot *testing.T) {
	t := MustFromMysqlFormatDefaultTZ("2001-01-01 00:00:00")
	kmgTest.Equal(t.Hour(), 0)
	kmgTest.Equal(t.Day(), 1)

	t = MustFromMysqlFormatDefaultTZ("0000-00-00 00:00:00")
	kmgTest.Equal(t.IsZero(), true)
}
