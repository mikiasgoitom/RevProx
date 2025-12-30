package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/RevProx/internal/contract"
)

type HealthHandler struct {
	healthCheckUseCase contract.IHealthCheckUseCase
	logger  contract.ILogger
}

func NewHealthCheckHandler(uc contract.IHealthCheckUseCase, l contract.ILogger) *HealthHandler {
	return &HealthHandler{healthCheckUseCase: uc, logger: l}
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	if err := h.healthCheckUseCase.Liveness(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	if err := h.healthCheckUseCase.Readiness(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
