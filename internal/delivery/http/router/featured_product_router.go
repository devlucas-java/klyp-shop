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
	adapter         *adapter.Adapter
}

func NewFeaturedProductRouter(
	jwt *jwt.JWTService,
	fh *handler.FeaturedProductHandler,
	log *logger.Logger,
	ur repository.UserRepository,
	a *adapter.Adapter,
) *FeaturedProductRouter {
	return &FeaturedProductRouter{
		jwtService:      jwt,
		featuredHandler: fh,
		log:             log,
		userRepository:  ur,
		adapter:         a,
	}
}

func (f *FeaturedProductRouter) RegisterFeaturedRoutes(r chi.Router) {

	r.Get("/", f.adapter.Adapt(f.featuredHandler.GetAllFeatured))
	r.Get("/seller/{sellerID}", f.adapter.Adapt(f.featuredHandler.GetFeaturedBySeller))

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(f.jwtService, f.log, f.userRepository))
		protected.Use(middleware.RoleMiddleware([]enums.Role{enums.SELLER, enums.ADMIN}, f.log))

		protected.Get("/me", f.adapter.Adapt(f.featuredHandler.GetMyFeatured))
		protected.Post("/", f.adapter.Adapt(f.featuredHandler.AddFeatured))
		protected.Delete("/{productID}", f.adapter.Adapt(f.featuredHandler.RemoveFeatured))
		protected.Patch("/{productID}/position", f.adapter.Adapt(f.featuredHandler.UpdatePosition))
	})
}
