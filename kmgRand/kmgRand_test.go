package kmgRand

import (
	"testing"

	. "github.com/bronze1man/kmg/kmgTest"
)

func TestIntBetween(t *testing.T) {
	Equal(IntBetween(0, 0), 0)
	hasView := [2]bool{}
	for i := 0; i < 100; i++ {
		ret := IntBetween(0, 1)
		Ok(ret == 0 || ret == 1)
		hasView[ret] = true
	}
	Equal(hasView[0], true)
	Equal(hasView[1], true)
}
