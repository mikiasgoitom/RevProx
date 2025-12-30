package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler exposes the Prometheus metrics endpoint.
type PrometheusHandler struct{}

// NewPrometheusHandler creates a new handler for Prometheus metrics.
func NewPrometheusHandler() *PrometheusHandler {
	return &PrometheusHandler{}
}

func (h *PrometheusHandler) GetMetrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}