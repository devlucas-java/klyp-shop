package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitAuthModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {
	userRepo := database.NewUserDB(db)
	userMapper := mapper.NewUserMapper()

	authService := service.NewAuthService(userRepo, jwtService, userMapper)

	authHandler := handler.NewAuthHandler(authService, log)
	adapter := adapter.NewAdapter(log)
	authRouter := router.NewAuthRouter(authHandler, jwtService, log, userRepo, adapter)

	r := chi.NewRouter()
	authRouter.RegisterAuthRoutes(r)

	return r
}
