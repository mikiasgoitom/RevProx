package usecase

import (
	"context"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/mikiasgoitom/RevProx/internal/contract"
	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
	valueobject "github.com/mikiasgoitom/RevProx/internal/domain/valueObject"
)

type ServeProxyRequestUseCase struct {
	ServeProxyRequestUseCase contract.IServeProxyRequestUseCase
	TimeService contract.ITimeService
	CacheRepository contract.ICacheRepository
	MetricsUseCase contract.IMetricsUseCase
	Logger contract.ILogger
	OriginRepository contract.IOriginRepository
	PolicyEvaluator contract.IPolicyEvaluator
	CachePolicy entity.CachePolicy
}

func NewServeProxyRequestUseCase(serveProxyRequestUseCase contract.IServeProxyRequestUseCase, timeService contract.ITimeService, cacheRepository contract.ICacheRepository, metricsUseCase contract.IMetricsUseCase, logger contract.ILogger, originRepository contract.IOriginRepository, PolicyEvaluator contract.IPolicyEvaluator, cachePolicy entity.CachePolicy) *ServeProxyRequestUseCase {
	return &ServeProxyRequestUseCase{
		ServeProxyRequestUseCase: serveProxyRequestUseCase,
		TimeService: timeService,
		CacheRepository: cacheRepository,
		MetricsUseCase: metricsUseCase,
		Logger: logger,
		OriginRepository: originRepository,
		PolicyEvaluator: PolicyEvaluator,
		CachePolicy: cachePolicy,
	}
}
func normalizeURL(req_url *url.URL) string {
	if req_url == nil {
		return ""
	}
	u := *req_url
	u.Fragment = ""
	
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if port != "" {
		if (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
			u.Host = host
		} else {
			u.Host = host + ":" + port
		}
	}

	if u.Path == "" {
		u.Path = "/"
	} else {
		u.Path = path.Clean(u.Path)
		if !strings.HasPrefix(u.Path, "/") {
			u.Path = "/" + u.Path
		}
	}
	queryParams := u.Query()
	if len(queryParams) == 0 {
		u.RawQuery = ""
		return u.String()
	}
	keys := make([]string, 0, len(queryParams))
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := url.Values{}
	for _, k := range keys {
		vals := queryParams[k]
		sort.Strings(vals)
		for _, v := range vals {
			out.Add(k, v)
		}
	}
	u.RawQuery = out.Encode()
	return u.String()
}
func (uc *ServeProxyRequestUseCase) ServeProxyRequest(ctx context.Context, req entity.RequestModel) (entity.ResponseModel, error) {
	// 1. Start timing.
	startTime := uc.TimeService.NowUnix()
	
    // 2. Normalize URL (sort query params, drop fragment) -> normalizedURL.
	normalizedURL := normalizeURL(req.URL)

    // 3. Build CacheKey {Method, NormalizedURL}.
	cacheKey := valueobject.CacheKey{
		Method: req.Method,
		NormalizedURL: normalizedURL,
	}

    // 4. Cache lookup, record cache latency; inc hit/miss metrics.
	cacheValRetrieved, found, err := uc.CacheRepository.Get(ctx, cacheKey)
	
	if err != nil {
		uc.Logger.Error(ctx, "Cache Get error", valueobject.LogField{Key: "error", Value: err.Error()}, valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL})
		return entity.ResponseModel{}, err
	}

	cacheLatency := uc.TimeService.NowUnix() - startTime
	if found {
		err = uc.MetricsUseCase.IncHit(ctx)
		if err != nil {
			uc.Logger.Error(ctx, "Metrics IncHit error", valueobject.LogField{Key: "error", Value: err.Error()})
		}
		err = uc.MetricsUseCase.RecordCacheLatency(ctx, cacheLatency)
		if err != nil {
			uc.Logger.Error(ctx, "Metrics RecordCacheLatency error", valueobject.LogField{Key: "error", Value: err.Error()})
		}
		uc.Logger.Info(ctx, "Cache hit", valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL}, valueobject.LogField{Key: "latency_ms", Value: cacheLatency})
		
		resp := cacheValRetrieved.Payload
		uc.Logger.Info(ctx, "Response served from cache", valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL}, valueobject.LogField{Key: "latency_ms", Value: cacheLatency})
		return resp, nil		
	} else {
		err = uc.MetricsUseCase.IncMiss(ctx)
		if err != nil {
			uc.Logger.Error(ctx, "Metrics IncMiss error", valueobject.LogField{Key: "error", Value: err.Error()})
		}
		err = uc.MetricsUseCase.RecordCacheLatency(ctx, cacheLatency)
		if err != nil {
			uc.Logger.Error(ctx, "Metrics RecordCacheLatency error", valueobject.LogField{Key: "error", Value: err.Error()})
		}
		uc.Logger.Info(ctx, "Cache miss", valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL}, valueobject.LogField{Key: "latency_ms", Value: cacheLatency})

	}

    // 5. Prepare origin request (preserve headers; add X-Forwarded-For, X-Forwarded-Host, X-Forwarded-Proto).
	originHeaders := req.Headers.Clone()
	if req.ClientIP != "" {
		originHeaders.Del("X-Forwarded-For")
		originHeaders.Add("X-Forwarded-For", req.ClientIP)
	}

	if req.URL != nil {
		originHeaders.Del("X-Forwarded-Host")
		originHeaders.Add("X-Forwarded-Host", req.URL.Host)
		originHeaders.Del("X-Forwarded-Proto")
		originHeaders.Add("X-Forwarded-Proto", req.URL.Scheme)
	}

    // 6. Origin.Fetch, record upstream latency
	originFetchStartTime := uc.TimeService.NowUnix()
	originReq := req
	originReq.Headers = originHeaders
	resp, err := uc.OriginRepository.Fetch(ctx, originReq)
	if err != nil {
		uc.Logger.Fatal(ctx, "Origin Fetch error", valueobject.LogField{Key: "error", Value: err.Error()}, valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL})
		return entity.ResponseModel{}, err
	}
	originFetchLatency := uc.TimeService.NowUnix() - originFetchStartTime
	err = uc.MetricsUseCase.RecordUpstreamLatency(ctx, originFetchLatency)
	if err != nil {
		uc.Logger.Error(ctx, "Metrics RecordUpstreamLatency error", valueobject.LogField{Key: "error", Value: err.Error()})
	}
	uc.Logger.Info(ctx, "Origin fetch successful", valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL}, valueobject.LogField{Key: "latency_ms", Value: originFetchLatency})

    // 7. Evaluate cacheability
	cacheable, ttl := uc.PolicyEvaluator.Evaluate(resp, req, uc.CachePolicy)
	resp.Cacheable = cacheable
	uc.Logger.Info(ctx, "Cache policy evaluated", valueobject.LogField{Key: "url", Value: normalizedURL}, valueobject.LogField{Key: "cacheable", Value: cacheable}, valueobject.LogField{Key: "ttl_seconds", Value: time.Unix(ttl, 0)})
    
	// 8. If cacheable and ttlSeconds > 0: build CacheEntry then Cache.Set(ctx, entry)
	if cacheable && ttl > 0{
		newCacheEntry := entity.CacheEntry{
			Key: cacheKey,
			Payload: resp,
			ExpiresAt: ttl,
			StoredAt: uc.TimeService.NowUnix(),
		}
		if err = uc.CacheRepository.Set(ctx, newCacheEntry.Key, newCacheEntry); err != nil {
			uc.Logger.Error(ctx, "Cache Set error", valueobject.LogField{Key: "error", Value: err.Error()}, valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL} )
		}
		uc.Logger.Info(ctx, "Response cached", valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL}, valueobject.LogField{Key: "ttl_seconds", Value: time.Unix(ttl, 0)} )
	}

    // 9. Update total latency metrics.
	totalLatency := uc.TimeService.NowUnix() - startTime
	err = uc.MetricsUseCase.RecordTotalLatency(ctx, totalLatency)
	if err != nil {
		uc.Logger.Error(ctx, "Metrics RecordTotalLatency error", valueobject.LogField{Key: "error", Value: err.Error()})
	}

    // 10. Log summary
	uc.Logger.Info(ctx, "Request served from origin", valueobject.LogField{Key: "method", Value: req.Method}, valueobject.LogField{Key: "url", Value: normalizedURL}, valueobject.LogField{Key: "cacheable", Value: cacheable}, valueobject.LogField{Key: "total_latency_ms", Value: totalLatency})

    // 11. Return ResponseModel
	return resp, nil

}