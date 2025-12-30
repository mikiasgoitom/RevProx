package handler

import "github.com/gin-gonic/gin"

type Router struct {
	healthCheckHandler *HealthHandler
	prometheusHandler  *PrometheusHandler
	proxyHandler       *ProxyHandler
}

func NewRouter(
	healthCheckHandler *HealthHandler,
	prometheusHandler *PrometheusHandler,
	proxyHandler *ProxyHandler,
) *Router {
	return &Router{
		healthCheckHandler: healthCheckHandler,
		prometheusHandler:  prometheusHandler,
		proxyHandler:       proxyHandler,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	baseUrl := router.Group("/api/v1")
	metrics := baseUrl.Group("/metrics")
	{
		metrics.GET("/prometheus", r.prometheusHandler.GetMetrics)
	}
	health := baseUrl.Group("/health")
	{
		health.GET("/livez", r.healthCheckHandler.Liveness)
		health.GET("/readyz", r.healthCheckHandler.Readiness)
	}
	proxy := baseUrl.Group("/proxy")
	{
        // This is the correct implementation for a catch-all proxy route.
        // "Any" matches all HTTP methods (GET, POST, PUT, etc.).
        // "/*path" is a wildcard that matches any path after /proxy/.
        proxy.Any("/*path", r.proxyHandler.HandleProxy)
    }
}