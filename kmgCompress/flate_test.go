package kmgCompress

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestFlate(ot *testing.T) {
	origin := []byte("")
	ob := FlateMustCompress([]byte(""))
	output := FlateMustUnCompress(ob)
	kmgTest.Equal(origin, output)
}
