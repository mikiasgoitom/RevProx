package contract

import (
	"context"

	"github.com/mikiasgoitom/RevProx/internal/domain/entity"

)

type IOriginRepository interface {
	Fetch(ctx context.Context, req entity.RequestModel) (entity.ResponseModel, error)
	HealthCheck(ctx context.Context) error
}
