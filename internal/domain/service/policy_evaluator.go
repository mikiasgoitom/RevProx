package domainservice

import (
	"log"
	"net/http"
	"time"

	"github.com/mikiasgoitom/RevProx/internal/contract"
	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
	"github.com/mikiasgoitom/RevProx/pkg/cachecontrol"
)

type PolicyEvaluator struct{}

func NewPolicyEvaluator() contract.IPolicyEvaluator {
	return &PolicyEvaluator{}
}

func (srv *PolicyEvaluator) Evaluate(resp entity.ResponseModel, req entity.RequestModel, cachePolicy entity.CachePolicy) (bool, int64) {
	if req.Method != http.MethodGet {
		return false, 0
	}
	cc := cachecontrol.Parse(resp.Headers.Get("cache-control"))

	if cachecontrol.Has(cc, "no-store") || cachecontrol.Has(cc, "private") {
		log.Println("Cache not allowed due to no-store or private directive")
		return false, 0
	}

	if !isCacheableStatusCode(resp.Status) {
		log.Printf("Status code %d is not cacheable\n", resp.Status)
		return false, 0
	}

	if sMaxAge, ok := cachecontrol.GetDuration(cc, "s-maxage"); ok {
		ttl := int64(time.Now().Unix() + int64(sMaxAge.Seconds()))
		log.Printf("Using s-maxage directive with duration %v and ttl %v\n", sMaxAge, ttl)
		return true, ttl
	}

	if maxAge, ok := cachecontrol.GetDuration(cc, "max-age"); ok {
		ttl := int64(time.Now().Unix() + int64(maxAge.Seconds()))
		log.Printf("Using max-age directive with duration %v and ttl %v\n", maxAge, ttl)
		return true, ttl
	}

	if expiresHeader := resp.Headers.Get("Expires"); expiresHeader != "" {
		expireTime, err := http.ParseTime(expiresHeader)
		if err == nil {
			ttl := time.Until(expireTime)
			if ttl > 0 {
				log.Printf("Using Expires header with ttl %v\n", ttl)
				return true, int64(ttl.Seconds())
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
