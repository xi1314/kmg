package kmgTime

import (
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestDurationFormat(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	t.Equal(DurationFormat(time.Second), "1s")
	t.Equal(DurationFormat(1000*time.Second), "16m40s")
	t.Equal(DurationFormat(1234*time.Millisecond), "1.23s")
	t.Equal(DurationFormat(1234*time.Microsecond), "1.23ms")
	t.Equal(DurationFormat(1234*time.Nanosecond), "1.23us")
}
