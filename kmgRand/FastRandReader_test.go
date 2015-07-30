package kmgRand

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestFastRandReader(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	buf := make([]byte, 10*1024)
	n, err := FastRandReader.Read(buf)
	t.Equal(n, 10*1024)
	t.Equal(err, nil)
}

func BenchmarkFastRandReader(ot *testing.B) {
	ot.StopTimer()
	t := kmgTest.NewTestTools(ot)
	buf := make([]byte, ot.N)
	ot.SetBytes(int64(1))
	ot.StartTimer()
	n, err := FastRandReader.Read(buf)
	t.Equal(n, ot.N)
	t.Equal(err, nil)
}
