package kmgEmail

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestIsEmailFormat(ot *testing.T) {
	kmgTest.Equal(IsEmailFormat("abc@gmail.com"), true)
	kmgTest.Equal(IsEmailFormat("abc"), false)
	kmgTest.Equal(IsEmailFormat("abc@asdf@gmail.com"), false)
	kmgTest.Equal(IsEmailFormat("@gmail.com"), false)
	kmgTest.Equal(IsEmailFormat("abc@"), false)
	kmgTest.Equal(IsEmailFormat("abc@abc"), false)
}
