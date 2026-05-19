package main

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/configs"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/internal/module"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	log := logger.NewLogger(logger.DEBUG)

	cfg := configs.InitConfig(log)
	db := configs.InitDB(log)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	jwtService := jwt.NewJWTService(cfg.JwtSecret)

	r.Mount("/api/v1/auth", module.InitAuthModule(db, log, jwtService))
	r.Mount("/api/v1/user", module.InitUserModule(db, log, jwtService))
	r.Mount("/api/v1/address", module.InitAddressModule(db, log, jwtService))
	r.Mount("/api/v1/seller", module.InitSellerModule(db, log, jwtService))
	r.Mount("/api/v1/product", module.InitProductModule(db, log, jwtService))
	r.Mount("/api/v1/order", module.InitOrderModule(db, log, jwtService))
	r.Mount("/api/v1/cart", module.InitShoppingCartModule(db, log, jwtService))
	r.Mount("/api/v1/dashboard", module.InitDashboardModule(db, log, jwtService))

	r.Mount("/api/v1/payment", module.InitPaymentModule(
		db, log, jwtService,
		cfg.GetBTCPayBaseURL(),
		cfg.GetBTCPayStoreID(),
		cfg.GetBTCPayAPIKey(),
		cfg.GetBTCPayWebhookSecret(),
	))

	chatModule, _ := module.InitChatModule(db, log, jwtService)
	r.Mount("/api/v1/chat", chatModule)

	r.Mount("/api/v1/featured", module.InitFeaturedProductModule(db, log, jwtService))

	log.Infof("Server is running on port %s", cfg.WebServerPort)
	if err := http.ListenAndServe(":"+cfg.WebServerPort, r); err != nil {
		log.Errorf("http listen err: %v", err)
		panic(err)
	}
}
