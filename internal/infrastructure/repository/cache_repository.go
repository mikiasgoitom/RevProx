package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/dgraph-io/ristretto"
	"github.com/mikiasgoitom/RevProx/internal/config"
	"github.com/mikiasgoitom/RevProx/internal/contract"
	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
	"github.com/mikiasgoitom/RevProx/internal/domain/valueobject"
)

type CacheRepository struct {
	cache *ristretto.Cache
}

func NewCacheRepository(cfg config.Config) (contract.ICacheRepository, error) {
	maxCostBytes, err := datasize.ParseString(cfg.Cache.MaxCost)
	if err != nil {
		return nil, fmt.Errorf("invalid cache max_cost '%s': %w", cfg.Cache.MaxCost, err)
	}
	bufferItem := cfg.Cache.BufferItems
	if bufferItem <= 0 {
		bufferItem = 64 // default buffer items
	}
	ristrettoConfig := &ristretto.Config{
		NumCounters: cfg.Cache.NumCounters,
		MaxCost:     int64(maxCostBytes.Bytes()),
		BufferItems: bufferItem,
	}
	cache, err := ristretto.NewCache(ristrettoConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ristretto cache: %w", err)
	}
	return &CacheRepository{
		cache: cache,
	}, nil

}

func (r *CacheRepository) Get(ctx context.Context, key valueobject.CacheKey) (entity.CacheEntry, bool, error) {
	cacheKey := fmt.Sprintf("%s:%s", key.Method, key.NormalizedURL)
	value, found := r.cache.Get(cacheKey)
	if !found {
		return entity.CacheEntry{}, false, nil
	}
	entry, ok := value.(entity.CacheEntry)
	if !ok {
		return entity.CacheEntry{}, false, fmt.Errorf("failed to cast cache value to CacheEntry")
	}
	return entry, true, nil
}

func (r *CacheRepository) Set(ctx context.Context, entry entity.CacheEntry) error {
	ttl := time.Until(time.Unix(entry.ExpiresAt, 0))
	if ttl <= 0 {
		return nil // Do not cache expired entries
	}

	// the cost is normally the size in bytes of the value being cached
	// for simplicity, we use the length of the body here
	cost := int64(len(entry.Payload.Body))
	if cost == 0 {
		cost = 1 // minimum cost
	}
	cacheKey := fmt.Sprintf("%s:%s", entry.Key.Method, entry.Key.NormalizedURL)
	wasAdded := r.cache.SetWithTTL(cacheKey, entry, cost, ttl)

	if !wasAdded {
		return fmt.Errorf("failed to add entry to cache")
	}

	// wait for the value to pass through the buffer for it to be available for subsequent gets
	r.cache.Wait()
	return nil

}

func (r *CacheRepository) HealthCheck(ctx context.Context) error {
	dummyKey := "healthcheck:key"
	dummyValue := "healthcheck:value"

	// set a dummy value
	if ok := r.cache.Set(dummyKey, dummyValue, 1); !ok {
		return fmt.Errorf("cache health check failed: unable to set value")
	}
	r.cache.Wait()
	// get the dummy value
	_, found := r.cache.Get(dummyKey)
	if !found {
		return fmt.Errorf("cache health check failed: unable to get value")
	}
	return nil
}
