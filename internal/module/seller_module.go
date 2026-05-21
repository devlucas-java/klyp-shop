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

func InitSellerModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {

	userRepository := database.NewUserDB(db)
	sellerRepository := database.NewSellerDB(db)
	sellerMapper := mapper.NewSellerMapper()

	sellerService := service.NewSellerService(log, userRepository, sellerRepository, sellerMapper)
	sellerHandler := handler.NewSellerHandler(sellerService, log)
	adapter := adapter.NewAdapter(log)
	sellerRouter := router.NewSellerRouter(jwtService, sellerHandler, log, userRepository, sellerRepository, adapter)

	r := chi.NewRouter()
	sellerRouter.RegisterSellerRoutes(r)

	return r
}
