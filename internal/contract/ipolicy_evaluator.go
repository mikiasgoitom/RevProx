package contract

import "github.com/mikiasgoitom/RevProx/internal/domain/entity"

type IPolicyEvaluator interface {
	Evaluate(resp entity.ResponseModel, req entity.RequestModel, cachePolicy entity.CachePolicy) (bool, int64)
}
