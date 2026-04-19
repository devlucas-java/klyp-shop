package router

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/go-chi/chi"
)

type AuthRouter struct {
	handler    *handler.AuthHandler
	jwtService *jwt.JWTService
}

func NewAuthRouter(handler *handler.AuthHandler, jwtService *jwt.JWTService) *AuthRouter {
	return &AuthRouter{
		handler:    handler,
		jwtService: jwtService,
	}
}

func (a *AuthRouter) RegisterRoutes(r chi.Router) {
	r.Post("/login", a.handler.Login)
	r.Post("/register", a.handler.Register)

	r.Route("/", func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware(a.jwtService))
		protected.Put("/password", a.handler.ChangePassword)
		protected.Post("/password", a.handler.VerifyPassword)
	})
}
