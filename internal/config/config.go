package config

import (
	"time"

	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
	valueobject "github.com/mikiasgoitom/RevProx/internal/domain/valueObject"
)

type Config struct {
	Server ServerConfig
	Cache  CacheConfig
	Origin OriginConfig
}

type ServerConfig struct {
	Port       string `mapstructure:"port"`
	Production bool   `mapstructure:"production"`
}
type CacheConfig struct {
	MaxCost     string       `mapstructure:"max_cost"`
	NumCounters int64        `mapstructure:"num_counters"`
	BufferItems int64        `mapstructure:"buffer_items"`
	Policy      PolicyConfig `mapstructure:"policy"`
}
type OriginConfig struct {
	OriginUrl string `mapstructure:"origin_url"`
}

type PolicyConfig struct {
	DefaultTTLSeconds       int64 `mapstructure:"default_ttl_seconds"`
	RespectNoCache          bool  `mapstructure:"respect_no_cache"`
	RespectNoStore          bool  `mapstructure:"respect_no_store"`
	RevalidateWindowSeconds int64 `mapstructure:"revalidate_window_seconds"`
}

func (pc *CacheConfig) ToCachePolicyEntity() entity.CachePolicy {
	return entity.CachePolicy{
		DefaultTTL:       valueobject.TTL{Duration: time.Duration(pc.Policy.DefaultTTLSeconds) * time.Second},
		RespectNoCache:   pc.Policy.RespectNoCache,
		RespectNoStore:   pc.Policy.RespectNoStore,
		RevalidateWindow: time.Duration(pc.Policy.RevalidateWindowSeconds) * time.Second,
	}
}
