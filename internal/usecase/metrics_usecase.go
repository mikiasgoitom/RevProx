package usecase

import (
	"context"

	"github.com/mikiasgoitom/caching-proxy/internal/contract"
	"github.com/mikiasgoitom/caching-proxy/internal/domain/entity"
)

type MetricsUseCase struct{}

func NewMetricsUseCase() contract.IMetricsUseCase {
	return &MetricsUseCase{}
}

func (uc *MetricsUseCase) IncHit(ctx context.Context) error                               {}
func (uc *MetricsUseCase) IncMiss(ctx context.Context) error                              {}
func (uc *MetricsUseCase) RecordUpstreamLatency(ctx context.Context, latency int64) error {}
func (uc *MetricsUseCase) RecordCacheLatency(ctx context.Context, latency int64) error    {}
func (uc *MetricsUseCase) RecordTotalLatency(ctx context.Context, latency int64) error    {}

// IncEviction(ctx context.Context) error
func (uc *MetricsUseCase) Snapshot(ctx context.Context) (entity.Metrics, error) {}
func (uc *MetricsUseCase) Reset(ctx context.Context) error                      {}
