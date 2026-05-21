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

type OrderRouter struct {
	jwtService     *jwt.JWTService
	orderHandler   *handler.OrderHandler
	log            *logger.Logger
	userRepository repository.UserRepository
	adapter        *adapter.Adapter
}

func NewOrderRouter(
	jwt *jwt.JWTService,
	oh *handler.OrderHandler,
	log *logger.Logger,
	ur repository.UserRepository,
	a *adapter.Adapter) *OrderRouter {
	return &OrderRouter{
		jwtService:     jwt,
		orderHandler:   oh,
		log:            log,
		userRepository: ur,
		adapter:        a,
	}
}

func (r *OrderRouter) RegisterOrderRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(r.jwtService, r.log, r.userRepository))

		protected.Post("/", r.adapter.Adapt(r.orderHandler.CreateOrder))
		protected.Get("/", r.adapter.Adapt(r.orderHandler.ListOrders))
		protected.Get("/{id}", r.adapter.Adapt(r.orderHandler.GetOrderByID))
		protected.Delete("/{id}", r.adapter.Adapt(r.orderHandler.CancelOrder))
	})
}
