package entity

import "time"

type Metrics struct {
	Hits            uint64
	Misses          uint64
	UpstreamLatency time.Duration
	CacheLatency    time.Duration
	TotalLatency    time.Duration
	Evictions       uint64
}
