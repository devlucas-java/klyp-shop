package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitAuthModule(db *gorm.DB, log *logger.Logger, jwtSecret string) chi.Router {
	userRepo := database.NewUserRepository(db, log)
	jwtSvc := jwt.NewJWTService(jwtSecret)

	authService := service.NewAuthService(userRepo, jwtSvc)

	authHandler := handler.NewAuthHandler(authService, log)
	authRouter := router.NewAuthRouter(authHandler, jwtSvc)

	r := chi.NewRouter()
	authRouter.RegisterRoutes(r)

	return r
}
