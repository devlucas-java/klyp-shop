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

func InitShoppingCartModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {
	shoppingCartRepository := database.NewShoppingCartDB(db, log)
	productRepository := database.NewProductDB(db, log)
	userRepository := database.NewUserDB(db, log)

	shoppingCartService := service.NewShoppingCartService(log, shoppingCartRepository, mapper.NewShoppingCartMapper())
	shoppingCartItemService := service.NewShoppingCartItemService(log, shoppingCartRepository, productRepository, mapper.NewShoppingCartMapper())
	shoppingCartHandler := handler.NewShoppingCartHandler(shoppingCartService, log)
	shoppingCartItemHandler := handler.NewShoppingCartItemHandler(shoppingCartItemService, log)
	shoppingCartRouter := router.NewShoppingCartRouter(jwtService, shoppingCartHandler, log, userRepository)
	shoppingCartItemRouter := router.NewShoppingCartItemRouter(shoppingCartItemHandler)

	r := chi.NewRouter()
	shoppingCartRouter.RegisterShoppingCartRoutes(r)
	shoppingCartItemRouter.RegisterShoppingCartItemRoutes(r)
	return r
}
