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
}

func NewSellerRouter(
	jwtService *jwt.JWTService,
	sellerHandler *handler.SellerHandler,
	log *logger.Logger,
	userRepository repository.UserRepository,
	sellerRepository repository.SellerRepository,
) *SellerRouter {
	return &SellerRouter{
		jwtService:       jwtService,
		sellerHandler:    sellerHandler,
		log:              log,
		userRepository:   userRepository,
		sellerRepository: sellerRepository,
	}
}

func (s *SellerRouter) RegisterSellerRoutes(r chi.Router) {
	r.Use(middleware.AuthMiddleware(s.jwtService, s.log, s.userRepository))

	r.Post("/seller", adapter.Adapt(s.sellerHandler.CreateSeller))
	r.Patch("/seller", adapter.Adapt(s.sellerHandler.UpdateSeller))
	r.Delete("/seller", adapter.Adapt(s.sellerHandler.DeleteSeller))
	r.Get("/seller/{id}", adapter.Adapt(s.sellerHandler.GetSellerByID))
}
