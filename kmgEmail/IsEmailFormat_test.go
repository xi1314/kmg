package kmgEmail

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestIsEmailFormat(ot *testing.T) {
	kmgTest.Equal(IsEmailFormat("abc@gmail.com"), true)
	kmgTest.Equal(IsEmailFormat("abc"), false)
	kmgTest.Equal(IsEmailFormat("abc@asdf@gmail.com"), false)
	kmgTest.Equal(IsEmailFormat("@gmail.com"), false)
	kmgTest.Equal(IsEmailFormat("abc@"), false)
	kmgTest.Equal(IsEmailFormat("abc@abc"), false)
}
