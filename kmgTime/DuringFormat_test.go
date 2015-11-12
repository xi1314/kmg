package kmgTime

import (
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestDurationFormat(ot *testing.T) {
	kmgTest.Equal(DurationFormat(time.Second), "1s")
	kmgTest.Equal(DurationFormat(1000*time.Second), "16m40s")
	kmgTest.Equal(DurationFormat(1234*time.Millisecond), "1.23s")
	kmgTest.Equal(DurationFormat(1234*time.Microsecond), "1.23ms")
	kmgTest.Equal(DurationFormat(1234*time.Nanosecond), "1.23µs")
}
