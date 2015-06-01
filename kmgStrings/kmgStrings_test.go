package kmgStrings

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestIsAllAphphabet(ot *testing.T) {
	kmgTest.Equal(IsAllAphphabet("abc"), true)
	kmgTest.Equal(IsAllAphphabet(""), true)
	kmgTest.Equal(IsAllAphphabet("123"), false)
}
