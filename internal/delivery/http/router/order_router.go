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
}

func NewOrderRouter(jwtService *jwt.JWTService, orderHandler *handler.OrderHandler, log *logger.Logger, userRepository repository.UserRepository) *OrderRouter {
	return &OrderRouter{
		jwtService:     jwtService,
		orderHandler:   orderHandler,
		log:            log,
		userRepository: userRepository,
	}
}

func (r *OrderRouter) RegisterOrderRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(r.jwtService, r.log, r.userRepository))

		protected.Post("/", adapter.Adapt(r.orderHandler.CreateOrder))
		protected.Get("/", adapter.Adapt(r.orderHandler.ListOrders))
		protected.Get("/{id}", adapter.Adapt(r.orderHandler.GetOrderByID))
		protected.Delete("/{id}", adapter.Adapt(r.orderHandler.CancelOrder))
	})
}
