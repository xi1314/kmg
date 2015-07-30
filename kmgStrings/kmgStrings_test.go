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

func TestFirstLetterToUpper(t *testing.T) {
	s := FirstLetterToUpper("title")
	kmgTest.Equal("Title", s)

	s1 := FirstLetterToUpper("")
	kmgTest.Equal("", s1)

	s2 := FirstLetterToUpper("123")
	kmgTest.Equal("123", s2)

	s3 := FirstLetterToUpper("中文")
	kmgTest.Equal("中文", s3)
}
