package utils

import "time"

type Clock interface {
	Now() time.Time
}

// RealClock implements the Clock interface using the real time.Now()
type RealClock struct{}

func (c *RealClock) Now() time.Time {
	return time.Now()
}
