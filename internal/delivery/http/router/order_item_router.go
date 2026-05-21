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
	adapter          *adapter.Adapter
}

func NewOrderItemRouter(
	oh *handler.OrderItemHandler,
	jwt *jwt.JWTService,
	log *logger.Logger,
	ur repository.UserRepository,
	a *adapter.Adapter) *OrderItemRouter {
	return &OrderItemRouter{
		orderItemHandler: oh,
		jwtService:       jwt,
		log:              log,
		userRepository:   ur,
		adapter:          a,
	}
}

func (r *OrderItemRouter) RegisterOrderItemRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(r.jwtService, r.log, r.userRepository))

		protected.Get("/{id}/items", r.adapter.Adapt(r.orderItemHandler.GetOrderItems))
		protected.Get("/{id}/items/{itemId}", r.adapter.Adapt(r.orderItemHandler.GetOrderItem))
	})
}
