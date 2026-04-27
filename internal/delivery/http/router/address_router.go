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

type AddressRouter struct {
	handler        *handler.AddressHandler
	log            *logger.Logger
	jwtService     *jwt.JWTService
	userRepository repository.UserRepository
}

func NewAddressRouter(h *handler.AddressHandler, l *logger.Logger, js *jwt.JWTService, ur repository.UserRepository) *AddressRouter {
	return &AddressRouter{
		handler:        h,
		log:            l,
		jwtService:     js,
		userRepository: ur,
	}
}

func (a *AddressRouter) Handle(router chi.Router) {

	router.Route("/", func(protect chi.Router) {
		protect.Use(middleware.AuthMiddleware(a.jwtService, a.log, a.userRepository))

		protect.Get("/", adapter.Adapt(a.handler.GetAddresses))
		protect.Post("/", adapter.Adapt(a.handler.CreateAddress))
		protect.Put("/{id}", adapter.Adapt(a.handler.UpdateAddress))
		protect.Delete("/{id}", adapter.Adapt(a.handler.DeleteAddress))
	})

}
