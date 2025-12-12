package entity

import valueobject "github.com/mikiasgoitom/RevProx/internal/domain/valueObject"

type CacheEntry struct {
	Key       valueobject.CacheKey
	Payload   ResponseModel
	ExpiresAt int64
	StoredAt  int64
}
