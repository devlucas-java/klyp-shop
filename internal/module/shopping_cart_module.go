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

func InitShoppingCartModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {
	shoppingCartRepository := database.NewShoppingCartDB(db)
	productRepository := database.NewProductDB(db)
	userRepository := database.NewUserDB(db)

	shoppingCartItemRepository := database.NewShoppingCartItemDB(db)
	shoppingCartService := service.NewShoppingCartService(log, shoppingCartRepository, mapper.NewShoppingCartMapper())
	shoppingCartItemService := service.NewShoppingCartItemService(log, shoppingCartRepository, shoppingCartItemRepository, productRepository, mapper.NewShoppingCartMapper())
	shoppingCartHandler := handler.NewShoppingCartHandler(shoppingCartService, log)
	shoppingCartItemHandler := handler.NewShoppingCartItemHandler(shoppingCartItemService, log)

	adapter := adapter.NewAdapter(log)
	shoppingCartRouter := router.NewShoppingCartRouter(jwtService, shoppingCartHandler, log, userRepository, adapter)
	shoppingCartItemRouter := router.NewShoppingCartItemRouter(jwtService, shoppingCartItemHandler, log, userRepository, adapter)

	r := chi.NewRouter()
	shoppingCartRouter.RegisterShoppingCartRoutes(r)
	shoppingCartItemRouter.RegisterShoppingCartItemRoutes(r)
	return r
}
