package contract

import (
	"context"

	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
)

type IProxyUseCase interface {
	ServeProxyRequest(ctx context.Context, req entity.RequestModel) (entity.ResponseModel, error)
}
