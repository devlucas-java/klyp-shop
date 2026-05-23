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

type AddressRouter struct {
	handler        *handler.AddressHandler
	adapter        *adapter.Adapter
	jwtService     *jwt.JWTService
	userRepository repository.UserRepository
	log            *logger.Logger
}

func NewAddressRouter(
	h *handler.AddressHandler,
	js *jwt.JWTService,
	ur repository.UserRepository,
	l *logger.Logger,
	a *adapter.Adapter,
) *AddressRouter {
	return &AddressRouter{
		handler:        h,
		adapter:        a,
		jwtService:     js,
		userRepository: ur,
		log:            l,
	}
}

func (a *AddressRouter) Handle(mux chi.Router) {
	mux.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(a.jwtService, a.log, a.userRepository))
		protected.Use(middleware.RoleMiddleware([]enums.Role{enums.USER}, a.log))

		protected.Get("/", a.adapter.Adapt(a.handler.GetAddresses))
		protected.Post("/", a.adapter.Adapt(a.handler.CreateAddress))
		protected.Put("/{id}", a.adapter.Adapt(a.handler.UpdateAddress))
		protected.Delete("/{id}", a.adapter.Adapt(a.handler.DeleteAddress))
	})
}
