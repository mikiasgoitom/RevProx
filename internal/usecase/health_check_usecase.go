package usecase

import (
  "context"

  "github.com/mikiasgoitom/RevProx/internal/contract"
  valueobject "github.com/mikiasgoitom/RevProx/internal/domain/valueObject"
)

type HealthCheckUseCase struct {
  OriginRepository contract.IOriginRepository
  CacheRepository contract.ICacheRepository
  Logger contract.ILogger
}

func NewHealthCheckUseCase(Logger contract.ILogger, orignrepo contract.IOriginRepository, cacherepo contract.ICacheRepository) contract.IHealthCheckUseCase {
  return &HealthCheckUseCase{
    Logger: Logger,
    OriginRepository: orignrepo,
    CacheRepository: cacherepo,
  }
}

func (uc *HealthCheckUseCase) Readyness(ctx context.Context) error {
  if err:= uc.OriginRepository.HealthCheck(ctx); err != nil {
    uc.Logger.Error(ctx, "Origin repository health check failed: %v", valueobject.LogField{Key: "Error: ", Value: err.Error()})
    return err
  }
  if err:= uc.CacheRepository.HealthCheck(ctx); err != nil {
    uc.Logger.Error(ctx, "Cache repository health check failed: %v", valueobject.LogField{Key: "Error: ", Value: err.Error()})
    return err
  }
  uc.Logger.Info(ctx, "Readyness check passed")
  return nil
}
func (uc *HealthCheckUseCase) Liveness(ctx context.Context) error {
  uc.Logger.Info(ctx, "Liveness check passed")
  return nil
}