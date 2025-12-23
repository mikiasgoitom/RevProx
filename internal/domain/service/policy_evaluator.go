package domainservice

import (
	"net/http"
	"time"

	"github.com/mikiasgoitom/caching-proxy/internal/contract"
	"github.com/mikiasgoitom/caching-proxy/internal/domain/entity"
	"github.com/mikiasgoitom/caching-proxy/pkg/cachecontrol"
)

type PolicyEvaluator struct{}

func NewPolicyEvaluator() contract.IPolicyEvaluator {
	return &PolicyEvaluator{}
}

func (srv *PolicyEvaluator) Evaluate(resp entity.ResponseModel, req entity.RequestModel, cachePolicy entity.CachePolicy) (bool, int64) {
	if req.Method != http.MethodGet {
		return false, 0
	}
	cc := cachecontrol.Parse(resp.Header.Get("cache-control"))

	if cachecontrol.Has(cc, "no-store") || cachecontrol.Has(cc, "private") {
		return false, 0
	}

	if !isCacheableStatusCode(resp.Status) {
		return false, 0
	}

	if sMaxAge, ok := cachecontrol.GetDuration(cc, "s-maxage"); ok {
		return true, int64(sMaxAge.Seconds())
	}

	if maxAge, ok := cachecontrol.GetDuration(cc, "max-age"); ok {
		return true, int64(maxAge.Seconds())
	}

	if expiresHeader := resp.Header.Get("Expires"); expiresHeader != "" {
		expireTime, err := http.ParseTime(expiresHeader)
		if err == nil {
			ttl := time.Until(expireTime)
			if ttl > 0 {
				return true, int64(ttl)
			} // no else because we can set out owm ttl
		}
	}

	if cachePolicy.DefaultTTL.Duration > 0 {
		return true, int64(cachePolicy.DefaultTTL.Duration.Seconds())
	}

	return false, 0
}

func isCacheableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusOK, // 200
		http.StatusNonAuthoritativeInfo, // 203
		http.StatusNoContent,            // 204
		http.StatusPartialContent,       // 206
		http.StatusMultipleChoices,      // 300
		http.StatusMovedPermanently,     // 301
		http.StatusNotFound,             // 404
		http.StatusMethodNotAllowed,     // 405
		http.StatusGone,                 // 410
		http.StatusNotImplemented:       // 501
		return true
	default:
		return false
	}
}
