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
	adapter        *adapter.Adapter
}

func NewUserRouter(
	jwt *jwt.JWTService,
	uh *handler.UserHandler,
	log *logger.Logger,
	ur repository.UserRepository,
	a *adapter.Adapter,
) *UserRouter {
	return &UserRouter{
		jwtService:     jwt,
		userHandler:    uh,
		log:            log,
		userRepository: ur,
		adapter:        a,
	}
}

func (u *UserRouter) RegisterUserRoutes(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(u.jwtService, u.log, u.userRepository))

		protected.Get("/me", u.adapter.Adapt(u.userHandler.GetMe))
		protected.Delete("/me", u.adapter.Adapt(u.userHandler.DeleteMe))
		protected.Patch("/me", u.adapter.Adapt(u.userHandler.UpdateMe))

		protected.Group(func(admin chi.Router) {
			admin.Use(middleware.RoleMiddleware([]enums.Role{enums.ADMIN}))

			admin.Post("/promote/{id}", u.adapter.Adapt(u.userHandler.PromoteUser))
			admin.Post("/demote/{id}", u.adapter.Adapt(u.userHandler.DemoteUser))
		})
	})
}
