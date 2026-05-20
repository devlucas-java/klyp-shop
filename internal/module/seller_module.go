package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitSellerModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {

	userRepository := database.NewUserDB(db, log)
	sellerRepository := database.NewSellerDB(db, log)
	sellerMapper := mapper.NewSellerMapper()

	sellerService := service.NewSellerService(log, userRepository, sellerRepository, sellerMapper)
	sellerHandler := handler.NewSellerHandler(sellerService, log)
	sellerRouter := router.NewSellerRouter(jwtService, sellerHandler, log, userRepository, sellerRepository)

	r := chi.NewRouter()
	sellerRouter.RegisterSellerRoutes(r)

	return r
}
