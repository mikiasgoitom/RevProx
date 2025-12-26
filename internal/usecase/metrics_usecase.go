package usecase

import (
  "context"

  "github.com/mikiasgoitom/RevProx/internal/contract"
  "github.com/mikiasgoitom/RevProx/internal/domain/entity"
)

type MetricsUseCase struct{}

// IncHit implements [contract.IMetricsUseCase].
func (m *MetricsUseCase) IncHit(ctx context.Context) error {
  panic("unimplemented")
}

// IncMiss implements [contract.IMetricsUseCase].
func (m *MetricsUseCase) IncMiss(ctx context.Context) error {
  panic("unimplemented")
}

// RecordCacheLatency implements [contract.IMetricsUseCase].
func (m *MetricsUseCase) RecordCacheLatency(ctx context.Context, latency int64) error {
  panic("unimplemented")
}

// RecordTotalLatency implements [contract.IMetricsUseCase].
func (m *MetricsUseCase) RecordTotalLatency(ctx context.Context, latency int64) error {
  panic("unimplemented")
}

// RecordUpstreamLatency implements [contract.IMetricsUseCase].
func (m *MetricsUseCase) RecordUpstreamLatency(ctx context.Context, latency int64) error {
  panic("unimplemented")
}

// Reset implements [contract.IMetricsUseCase].
func (m *MetricsUseCase) Reset(ctx context.Context) error {
  panic("unimplemented")
}

// Snapshot implements [contract.IMetricsUseCase].
func (m *MetricsUseCase) Snapshot(ctx context.Context) (entity.Metrics, error) {
  panic("unimplemented")
}

func NewMetricsUseCase() contract.IMetricsUseCase {
  return &MetricsUseCase{}
}
