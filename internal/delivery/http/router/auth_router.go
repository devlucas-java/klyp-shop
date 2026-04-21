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

type AuthRouter struct {
	handler        *handler.AuthHandler
	jwtService     *jwt.JWTService
	log            *logger.Logger
	userRepository repository.UserRepository
}

func NewAuthRouter(handler *handler.AuthHandler, jwtService *jwt.JWTService, log *logger.Logger, userRepository repository.UserRepository) *AuthRouter {
	return &AuthRouter{
		handler:        handler,
		jwtService:     jwtService,
		log:            log,
		userRepository: userRepository,
	}
}

func (a *AuthRouter) RegisterAuthRoutes(r chi.Router) {
	r.Post("/login", adapter.Adapt(a.handler.Login))
	r.Post("/register", adapter.Adapt(a.handler.Register))

	r.Route("/", func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware(a.jwtService, a.log, a.userRepository))
		protected.Put("/password", adapter.Adapt(a.handler.ChangePassword))
		protected.Post("/password", adapter.Adapt(a.handler.VerifyPassword))
	})
}
