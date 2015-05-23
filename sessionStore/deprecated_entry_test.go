package sessionStore

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func Test(ot *testing.T) {
	kmgTest.TestWarpper(ot, &Tester{})
}

type Tester struct {
	kmgTest.TestTools
}
