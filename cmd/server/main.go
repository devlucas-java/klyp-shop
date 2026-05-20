package main

import (
	"github.com/devlucas-java/klyp-shop/configs"
	httpDelivery "github.com/devlucas-java/klyp-shop/internal/delivery/http"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log := logger.NewLogger(logger.DEBUG)

	cfg := configs.InitConfig(log)
	db := configs.InitDB(log)

	reg := prometheus.NewRegistry()
	metric := metrics.NewMetric(reg)
	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	jwtService := jwt.NewJWTService(cfg.JwtSecret, cfg.JwtExpireIn)

	router := httpDelivery.NewRouter(httpDelivery.RouterDeps{
		Cfg:            cfg,
		DB:             db,
		Log:            log,
		JwtService:     jwtService,
		Metric:         metric,
		MetricsHandler: metricsHandler,
	})

	server := httpDelivery.NewServer(cfg.WebServerPort, router, log)
	if err := server.Run(); err != nil {
		log.Errorf("server error: %v", err)
		panic(err)
	}
}
