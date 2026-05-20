package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitOrderModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService, metric *metrics.Metric) chi.Router {
	orderRepository := database.NewOrderDB(db, log)
	orderItemRepository := database.NewOrderItemDB(db, log)
	productRepository := database.NewProductDB(db, log)
	userRepository := database.NewUserDB(db, log)
	addressRepository := database.NewAddressDB(db, log)

	orderService := service.NewOrderService(log, orderRepository, userRepository, addressRepository, productRepository, mapper.NewOrderMapper(), metric)
	orderItemService := service.NewOrderItemService(log, orderItemRepository, orderRepository, productRepository, mapper.NewOrderMapper())
	orderHandler := handler.NewOrderHandler(orderService, log)
	orderItemHandler := handler.NewOrderItemHandler(orderItemService, log)
	orderRouter := router.NewOrderRouter(jwtService, orderHandler, log, userRepository)
	orderItemRouter := router.NewOrderItemRouter(orderItemHandler, jwtService, log, userRepository)

	r := chi.NewRouter()
	orderRouter.RegisterOrderRoutes(r)
	orderItemRouter.RegisterOrderItemRoutes(r)
	return r
}
