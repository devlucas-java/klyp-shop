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

func (u *UserRouter) RegisterUserRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(u.jwtService, u.log, u.userRepository))

		protected.Get("/me", adapter.Adapt(u.userHandler.GetMe))
		protected.Delete("/me", adapter.Adapt(u.userHandler.DeleteMe))
		protected.Patch("/me", adapter.Adapt(u.userHandler.UpdateMe))

		protected.Group(func(admin chi.Router) {
			admin.Use(middleware.RoleMiddleware([]enums.Role{enums.ADMIN}))

			admin.Post("/promote/{id}", adapter.Adapt(u.userHandler.PromoteUser))
			admin.Post("/demote/{id}", adapter.Adapt(u.userHandler.DemoteUser))
		})
	})
}
