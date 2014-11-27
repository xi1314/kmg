package kmgGob

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func Test(t *testing.T) {
	kmgTest.TestWarpper(t, &S{})
}

type S struct {
	kmgTest.TestTools
}
