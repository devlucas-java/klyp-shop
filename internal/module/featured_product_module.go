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

func InitFeaturedProductModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {
	userRepository := database.NewUserDB(db, log)
	productRepository := database.NewProductDB(db)
	featuredRepository := database.NewFeaturedProductDB(db)

	featuredService := service.NewFeaturedProductService(log, featuredRepository, productRepository, userRepository)
	featuredHandler := handler.NewFeaturedProductHandler(featuredService, log)
	featuredRouter := router.NewFeaturedProductRouter(jwtService, featuredHandler, log, userRepository)

	r := chi.NewRouter()
	featuredRouter.RegisterFeaturedRoutes(r)
	return r
}
