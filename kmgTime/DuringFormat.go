package kmgTime

import (
	"time"
)

func DurationFormat(dur time.Duration) string {
	if dur >= time.Second {
		mod := 10 * time.Millisecond
		dur = (dur / mod) * mod
		return dur.String()
	}
	if dur >= time.Millisecond {
		mod := 10 * time.Microsecond
		dur = (dur / mod) * mod
		return dur.String()
	}
	if dur >= time.Microsecond {
		mod := 10 * time.Nanosecond
		dur = dur / mod * mod
		return dur.String()
	}
	return dur.String()
}
