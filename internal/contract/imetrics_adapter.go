package contract

import (
	"context"
	"time"
)

type IMetricsAdapter interface {
	IncHit(ctx context.Context) error
	IncMiss(ctx context.Context) error
	RecordEviction(ctx context.Context) error
	RecordUpstreamLatency(ctx context.Context, latency time.Duration) error
	RecordCacheLatency(ctx context.Context, latency time.Duration) error
	RecordTotalLatency(ctx context.Context, latency time.Duration) error
}
