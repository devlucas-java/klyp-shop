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

type ShoppingCartRouter struct {
	jwtService          *jwt.JWTService
	shoppingCartHandler *handler.ShoppingCartHandler
	log                 *logger.Logger
	userRepository      repository.UserRepository
}

func NewShoppingCartRouter(
	jwtService *jwt.JWTService,
	shoppingCartHandler *handler.ShoppingCartHandler,
	log *logger.Logger,
	userRepository repository.UserRepository,
) *ShoppingCartRouter {
	return &ShoppingCartRouter{
		jwtService:          jwtService,
		shoppingCartHandler: shoppingCartHandler,
		log:                 log,
		userRepository:      userRepository,
	}
}

func (r *ShoppingCartRouter) RegisterShoppingCartRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(r.jwtService, r.log, r.userRepository))

		protected.Get("/", adapter.Adapt(r.shoppingCartHandler.GetCart))
		protected.Delete("/", adapter.Adapt(r.shoppingCartHandler.ClearCart))
	})
}
