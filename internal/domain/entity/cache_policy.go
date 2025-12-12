package entity

import (
	"time"

	valueobject "github.com/mikiasgoitom/RevProx/internal/domain/valueObject"
)

type CachePolicy struct {
	DefaultTTL       valueobject.TTL
	RespectNoCache   bool
	RespectNoStore   bool
	RevalidateWindow time.Duration
}
