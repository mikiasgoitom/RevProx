package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/caching-proxy/internal/contract"
	"github.com/mikiasgoitom/caching-proxy/internal/domain/entity"
)

type OriginRepository struct {
	client      *http.Client
	originUrl   *url.URL
	timeService contract.ITimeService
}

func NewHttpOriginRepository(originUrl string, timeService contract.ITimeService) (contract.IOriginRepository, error) {
	parsedUrl, err := url.Parse(originUrl)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	return &OriginRepository{
		client:      &client,
		originUrl:   parsedUrl,
		timeService: timeService,
	}, nil
}

func (r *OriginRepository) Fetch(ctx context.Context, req entity.RequestModel) (entity.ResponseModel, error) {
	targetUrl := r.originUrl.ResolveReference(req.URL)
	originReq, err := http.NewRequestWithContext(ctx, req.Method, targetUrl.String(), bytes.NewReader(req.Body))
	if err != nil {
		return entity.ResponseModel{}, fmt.Errorf("failed to create origin request: %w", err)
	}
	originReq.Header = req.Headers.Clone()

	httpResp, err := r.client.Do(originReq)
	if err != nil {
		return entity.ResponseModel{}, fmt.Errorf("failed to perform origin request: %w", err)
	}
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return entity.ResponseModel{}, fmt.Errorf("failed to read origin response body: %w", err)
	}
	cacheControlHeader := httpResp.Header.Get("Cache-Control")
	response := entity.ResponseModel{
		ID:          uuid.New().String(),
		Status:      httpResp.StatusCode,
		Header:      httpResp.Header.Clone(),
		Body:        body,
		GeneratedAt: r.timeService.NowUnix(),
		Cacheable:   strings.Contains(cacheControlHeader, "public"),
	}
	return response, nil
}
func (r *OriginRepository) HealthCheck(ctx context.Context) error {
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(healthCtx, http.MethodHead, r.originUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}
	response, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return fmt.Errorf("origin service unhealthy, status code: %d", response.StatusCode)
	}
	return nil
}
