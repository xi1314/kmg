package kmgCompress

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestFlate(ot *testing.T) {
	origin := []byte("")
	ob := FlateMustCompress([]byte(""))
	output := FlateMustUnCompress(ob)
	kmgTest.Equal(origin, output)
}
