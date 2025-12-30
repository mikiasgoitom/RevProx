package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	domainservice "github.com/mikiasgoitom/RevProx/internal/domain/service"
	"github.com/mikiasgoitom/RevProx/internal/domain/valueobject"
	"github.com/mikiasgoitom/RevProx/internal/handler"
	"github.com/mikiasgoitom/RevProx/internal/infrastructure/configservice"
	"github.com/mikiasgoitom/RevProx/internal/infrastructure/logger"
	metricsadapter "github.com/mikiasgoitom/RevProx/internal/infrastructure/metrics_adapter"
	"github.com/mikiasgoitom/RevProx/internal/infrastructure/repository"
	"github.com/mikiasgoitom/RevProx/internal/infrastructure/timeservice"
	"github.com/mikiasgoitom/RevProx/internal/usecase"
)

func main() {

	// ---------------infrastructure implementation---------------
	cfgService := configservice.NewViperAdapter()
	cfg, err := cfgService.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	appLogger, err := logger.NewZapAdapter(cfg.Server.Production)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	appLogger.Info(context.Background(), "Configuration loaded successfully")
	timeService := timeservice.NewTimeService()
	originRepo, err := repository.NewHttpOriginRepository(cfg.Origin.OriginUrl, timeService)
	if err != nil {
		appLogger.Error(context.Background(), "failed to create origin repository", valueobject.LogField{Key: "error", Value: err})
		os.Exit(1)
	}
	cacheRepo, err := repository.NewCacheRepository(cfg)
	if err != nil {
		appLogger.Error(context.Background(), "failed to create cache repository", valueobject.LogField{Key: "error", Value: err})
		os.Exit(1)
	}
	prometheusMetrics := metricsadapter.NewPrometheusAdapter()
	policyEvaluator := domainservice.NewPolicyEvaluator()
	// ---------------usecase implementaion---------------

	proxyUsecase := usecase.NewProxyUsecase(timeService, cacheRepo, prometheusMetrics, appLogger, originRepo, policyEvaluator, cfg.Cache.ToCachePolicyEntity())
	healthCheckUsecase := usecase.NewHealthCheckUseCase(appLogger, originRepo, cacheRepo)

	// --------------- handler implementation---------------
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckUsecase, appLogger)
	prometheusHandler := handler.NewPrometheusHandler()
	proxyHandler := handler.NewProxyHandler(proxyUsecase, appLogger)

	// --------------- router setup---------------
	router := handler.NewRouter(healthCheckHandler, prometheusHandler, proxyHandler)

	ginEngine := gin.Default()

	router.SetupRoutes(ginEngine)

	// --------------- start server---------------
	appLogger.Info(context.Background(), "Starting server on port "+cfg.Server.Port)
	if err := ginEngine.Run(":" + cfg.Server.Port); err != nil {
		appLogger.Error(context.Background(), "failed to start server", valueobject.LogField{Key: "error", Value: err})
	}

}
