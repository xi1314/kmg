package kmgTime

import "time"

type Nower interface {
	Now() time.Time
}

type tDefaultNower struct{}

type FixedNower struct {
	Time time.Time
}

var DefaultNower tDefaultNower

func GetDefaultNower() Nower {
	return DefaultNower
}

func NewFixedNower(time time.Time) Nower {
	return FixedNower{time}
}

func (nower tDefaultNower) Now() time.Time {
	return time.Now()
}

func (nower FixedNower) Now() time.Time {
	return nower.Time
}
