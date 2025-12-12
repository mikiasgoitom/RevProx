package valueobject

import "time"

type TTL struct {
	Duration time.Duration
	Adaptive bool
	Min      time.Duration
	Max      time.Duration
}
