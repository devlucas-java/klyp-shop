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

type ProductRouter struct {
	jwtService     *jwt.JWTService
	productHandler *handler.ProductHandler
	log            *logger.Logger
	userRepository repository.UserRepository
	productService repository.ProductRepository
}

func NewProductRouter(jwtService *jwt.JWTService, productHandler *handler.ProductHandler, log *logger.Logger, userRepository repository.UserRepository, productService repository.ProductRepository) *ProductRouter {
	return &ProductRouter{
		jwtService:     jwtService,
		productHandler: productHandler,
		log:            log,
		userRepository: userRepository,
		productService: productService,
	}
}

func (p *ProductRouter) RegisterProductRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(p.jwtService, p.log, p.userRepository))

		protected.Post("/product", adapter.Adapt(p.productHandler.CreateProduct))
		protected.Get("/product/{id}", adapter.Adapt(p.productHandler.GetProductByID))
		protected.Patch("/product/{id}", adapter.Adapt(p.productHandler.UpdateProduct))
		protected.Delete("/product/{id}", adapter.Adapt(p.productHandler.DeleteProduct))
	})
}
