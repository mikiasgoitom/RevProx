package contract

import (
	"context"

	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
	"github.com/mikiasgoitom/RevProx/internal/domain/valueobject"
)

type ICacheRepository interface {
	Get(ctx context.Context, key valueobject.CacheKey) (entity.CacheEntry, bool, error)
	Set(ctx context.Context,value entity.CacheEntry) error
	HealthCheck(ctx context.Context) error
}
