package router

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type OrderItemRouter struct {
	orderItemHandler *handler.OrderItemHandler
	jwtService       *jwt.JWTService
	log              *logger.Logger
	userRepository   repository.UserRepository
}

func NewOrderItemRouter(orderItemHandler *handler.OrderItemHandler, jwtService *jwt.JWTService, log *logger.Logger, userRepository repository.UserRepository) *OrderItemRouter {
	return &OrderItemRouter{
		orderItemHandler: orderItemHandler,
		jwtService:       jwtService,
		log:              log,
		userRepository:   userRepository,
	}
}

func (r *OrderItemRouter) RegisterOrderItemRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(r.jwtService, r.log, r.userRepository))

		protected.Get("/{id}/items", adapter.Adapt(r.orderItemHandler.GetOrderItems))
		protected.Get("/{id}/items/{itemId}", adapter.Adapt(r.orderItemHandler.GetOrderItem))
	})
}
