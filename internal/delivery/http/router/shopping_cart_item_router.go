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

type ShoppingCartItemRouter struct {
	jwtService              *jwt.JWTService
	shoppingCartItemHandler *handler.ShoppingCartItemHandler
	log                     *logger.Logger
	userRepository          repository.UserRepository
	adapter                 *adapter.Adapter
}

func NewShoppingCartItemRouter(
	jwtService *jwt.JWTService,
	sh *handler.ShoppingCartItemHandler,
	log *logger.Logger,
	ur repository.UserRepository,
	adapter *adapter.Adapter,
) *ShoppingCartItemRouter {
	return &ShoppingCartItemRouter{
		jwtService:              jwtService,
		shoppingCartItemHandler: sh,
		log:                     log,
		userRepository:          ur,
		adapter:                 adapter,
	}
}

func (s *ShoppingCartItemRouter) RegisterShoppingCartItemRoutes(r chi.Router) {
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(s.jwtService, s.log, s.userRepository))

		protected.Post("/items", s.adapter.Adapt(s.shoppingCartItemHandler.AddItem))
		protected.Patch("/items/{id}", s.adapter.Adapt(s.shoppingCartItemHandler.UpdateItem))
		protected.Delete("/items/{id}", s.adapter.Adapt(s.shoppingCartItemHandler.RemoveItem))
	})
}
