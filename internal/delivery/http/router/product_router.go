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

func (p *ProductRouter) RegisterProductRoutes(protect chi.Router) {

	protect.Use(middleware.AuthMiddleware(p.jwtService, p.log, p.userRepository))

	protect.Post("/product", adapter.Adapt(p.productHandler.CreateProduct))
	protect.Get("/product/{id}", adapter.Adapt(p.productHandler.GetProductByID))
	protect.Patch("/product/{id}", adapter.Adapt(p.productHandler.UpdateProduct))
	protect.Delete("/product/{id}", adapter.Adapt(p.productHandler.DeleteProduct))
}
