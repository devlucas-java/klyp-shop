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

func InitProductModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {

	userRepository := database.NewUserDB(db, log)
	productRepository := database.NewProductDB(db)
	sellerRepository := database.NewSellerDB(db)
	productMapper := mapper.NewProductMapper()

	productService := service.NewProductService(log, productRepository, userRepository, sellerRepository, productMapper)
	productHandler := handler.NewProductHandler(productService, log)
	productRouter := router.NewProductRouter(jwtService, productHandler, log, userRepository, productRepository)

	r := chi.NewRouter()
	productRouter.RegisterProductRoutes(r)

	return r
}
