package entity

import (
	"time"

	"github.com/mikiasgoitom/RevProx/internal/domain/valueobject"
)

type CachePolicy struct {
	DefaultTTL valueobject.TTL
	RespectNoCache bool
	RespectNoStore bool
	RevalidateWindow time.Duration
}