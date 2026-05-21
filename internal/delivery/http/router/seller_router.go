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

type SellerRouter struct {
	jwtService       *jwt.JWTService
	sellerHandler    *handler.SellerHandler
	log              *logger.Logger
	userRepository   repository.UserRepository
	sellerRepository repository.SellerRepository
	adapter          *adapter.Adapter
}

func NewSellerRouter(
	jwt *jwt.JWTService,
	sh *handler.SellerHandler,
	log *logger.Logger,
	ur repository.UserRepository,
	sr repository.SellerRepository,
	a *adapter.Adapter,
) *SellerRouter {
	return &SellerRouter{
		jwtService:       jwt,
		sellerHandler:    sh,
		log:              log,
		userRepository:   ur,
		sellerRepository: sr,
		adapter:          a,
	}
}

func (s *SellerRouter) RegisterSellerRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(s.jwtService, s.log, s.userRepository))

		protected.Post("/seller", s.adapter.Adapt(s.sellerHandler.CreateSeller))
		protected.Patch("/seller", s.adapter.Adapt(s.sellerHandler.UpdateSeller))
		protected.Delete("/seller", s.adapter.Adapt(s.sellerHandler.DeleteSeller))
		protected.Get("/seller/{id}", s.adapter.Adapt(s.sellerHandler.GetSellerByID))
	})
}
