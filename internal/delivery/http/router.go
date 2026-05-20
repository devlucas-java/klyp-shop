package http

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/configs"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/internal/module"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
)

type RouterDeps struct {
	Cfg            *configs.Conf
	DB             *gorm.DB
	Log            *logger.Logger
	JwtService     *jwt.JWTService
	Metric         *metrics.Metric
	MetricsHandler http.Handler
}

func NewRouter(deps RouterDeps) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RecordMetricsMiddleware(deps.Metric))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)

	// ── Infra ─────────────────────────────────────────────────────────────
	r.Handle("/metrics", deps.MetricsHandler)
	r.Get("/health/check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// ── API v1 ────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/auth", module.InitAuthModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/user", module.InitUserModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/address", module.InitAddressModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/seller", module.InitSellerModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/product", module.InitProductModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/order", module.InitOrderModule(deps.DB, deps.Log, deps.JwtService, deps.Metric))
		api.Mount("/cart", module.InitShoppingCartModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/dashboard", module.InitDashboardModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/featured", module.InitFeaturedProductModule(deps.DB, deps.Log, deps.JwtService))
		api.Mount("/payment", module.InitPaymentModule(
			deps.DB, deps.Log, deps.JwtService,
			deps.Cfg.GetBTCPayBaseURL(),
			deps.Cfg.GetBTCPayStoreID(),
			deps.Cfg.GetBTCPayAPIKey(),
			deps.Cfg.GetBTCPayWebhookSecret(),
			deps.Metric,
		))

		chatModule, _ := module.InitChatModule(deps.DB, deps.Log, deps.JwtService, deps.Metric)
		api.Mount("/chat", chatModule)
	})

	return r
}
