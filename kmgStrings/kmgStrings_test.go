package kmgStrings

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestIsAllAphphabet(ot *testing.T) {
	kmgTest.Equal(IsAllAphphabet("abc"), true)
	kmgTest.Equal(IsAllAphphabet(""), true)
	kmgTest.Equal(IsAllAphphabet("123"), false)
}
