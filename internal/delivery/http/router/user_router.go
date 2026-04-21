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

type UserRouter struct {
	jwtService     *jwt.JWTService
	userHandler    *handler.UserHandler
	log            *logger.Logger
	userRepository repository.UserRepository
}

func NewUserRouter(jwtService *jwt.JWTService, userHandler *handler.UserHandler, log *logger.Logger, userRepository repository.UserRepository) *UserRouter {
	return &UserRouter{
		jwtService:     jwtService,
		userHandler:    userHandler,
		log:            log,
		userRepository: userRepository,
	}
}

func (u *UserRouter) RegisterUserRoutes(r chi.Router) {

	r.Route("/", func(protect chi.Router) {
		protect.Use(middleware.AuthMiddleware(u.jwtService, u.log, u.userRepository))
		protect.Get("/me", adapter.Adapt(u.userHandler.GetMe))
		protect.Delete("/me", adapter.Adapt(u.userHandler.DeleteMe))
		protect.Patch("/me", adapter.Adapt(u.userHandler.UpdateMe))
		protect.Post("/promote/{id}", adapter.Adapt(u.userHandler.PromoteUser))
		protect.Post("/demote/{id}", adapter.Adapt(u.userHandler.DemoteUser))

	})
}
