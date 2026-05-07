package router

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type FeaturedProductRouter struct {
	jwtService      *jwt.JWTService
	featuredHandler *handler.FeaturedProductHandler
	log             *logger.Logger
	userRepository  repository.UserRepository
}

func NewFeaturedProductRouter(
	jwtService *jwt.JWTService,
	featuredHandler *handler.FeaturedProductHandler,
	log *logger.Logger,
	userRepository repository.UserRepository,
) *FeaturedProductRouter {
	return &FeaturedProductRouter{
		jwtService:      jwtService,
		featuredHandler: featuredHandler,
		log:             log,
		userRepository:  userRepository,
	}
}

func (f *FeaturedProductRouter) RegisterFeaturedRoutes(r chi.Router) {
	// Public — anyone can see a seller's top 10
	r.Get("/seller/{sellerID}", adapter.Adapt(f.featuredHandler.GetFeaturedBySeller))

	// Protected — only sellers manage their own top 10
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware(f.jwtService, f.log, f.userRepository))
		protected.Use(middleware.RoleMiddleware([]enums.Role{enums.SELLER, enums.ADMIN}))

		protected.Get("/me", adapter.Adapt(f.featuredHandler.GetMyFeatured))
		protected.Post("/", adapter.Adapt(f.featuredHandler.AddFeatured))
		protected.Delete("/{productID}", adapter.Adapt(f.featuredHandler.RemoveFeatured))
		protected.Patch("/{productID}/position", adapter.Adapt(f.featuredHandler.UpdatePosition))
	})
}
