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
	adapter        *adapter.Adapter
}

func NewAuthRouter(
	h *handler.AuthHandler,
	jwt *jwt.JWTService,
	l *logger.Logger,
	ur repository.UserRepository,
	a *adapter.Adapter,
) *AuthRouter {
	return &AuthRouter{
		handler:        h,
		jwtService:     jwt,
		log:            l,
		userRepository: ur,
		adapter:        a,
	}
}

func (a *AuthRouter) RegisterAuthRoutes(r chi.Router) {
	r.Post("/login", a.adapter.Adapt(a.handler.Login))
	r.Post("/register", a.adapter.Adapt(a.handler.Register))

	r.Route("/", func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(a.jwtService, a.log, a.userRepository))
		protected.Put("/password", a.adapter.Adapt(a.handler.ChangePassword))
		protected.Post("/password", a.adapter.Adapt(a.handler.VerifyPassword))
	})
}
