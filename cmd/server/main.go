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

	cfg := configs.InitConfigDev(log)
	db := configs.InitDBDev(log)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	jwtService := jwt.NewJWTService(cfg.JwtSecret)

	authModule := module.InitAuthModule(db, log, jwtService)
	r.Mount("/api/v1/auth", authModule)

	userModule := module.InitUserModule(db, log, jwtService)
	r.Mount("/api/v1/user", userModule)

	log.Infof("Server is running on port %s", cfg.WebServerPort)
	if err := http.ListenAndServe(":"+cfg.WebServerPort, r); err != nil {
		log.Errorf("http listen err: %v", err)
		panic(err)
	}
}
