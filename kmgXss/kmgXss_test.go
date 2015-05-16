package kmgXss

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestUrlv(t *testing.T) {
	kmgTest.Equal(Urlv("abcd"), "abcd")
	kmgTest.Equal(Urlv("abcd哈"), "abcd%E5%93%88")
	kmgTest.Equal(Urlv("abcd "), "abcd%20")
	kmgTest.Equal(Urlv("abcd.abc"), "abcd.abc")

}

func TestH(t *testing.T) {
	kmgTest.Equal(H("<"), "&lt;")
	kmgTest.Equal(H("abc&"), "abc&amp;")
	kmgTest.Equal(H("abcdef"), "abcdef")
	kmgTest.Equal(H("abcd"), "abcd")
}
