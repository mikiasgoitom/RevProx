package contract

import (
	"context"

	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
)

type IMetricsUseCase interface {
	IncHit(ctx context.Context) error
	IncMiss(ctx context.Context) error
	RecordUpstreamLatency(ctx context.Context, latency int64) error
	RecordCacheLatency(ctx context.Context, latency int64) error
	RecordTotalLatency(ctx context.Context, latency int64) error
	// IncEviction(ctx context.Context) error
	Snapshot(ctx context.Context) (entity.Metrics, error)
	Reset(ctx context.Context) error
}
