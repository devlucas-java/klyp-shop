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

type PaymentRouter struct {
	jwtService     *jwt.JWTService
	handler        *handler.PaymentHandler
	log            *logger.Logger
	userRepository repository.UserRepository
	adapter        *adapter.Adapter
}

func NewPaymentRouter(
	jwt *jwt.JWTService,
	h *handler.PaymentHandler,
	log *logger.Logger,
	ur repository.UserRepository,
	a *adapter.Adapter,
) *PaymentRouter {
	return &PaymentRouter{
		jwtService:     jwt,
		handler:        h,
		log:            log,
		userRepository: ur,
		adapter:        a,
	}
}

func (p *PaymentRouter) RegisterPaymentRoutes(r chi.Router) {
	r.Post("/webhook", p.adapter.Adapt(p.handler.Webhook))

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(p.jwtService, p.log, p.userRepository))
		protected.Post("/orders/{orderID}/invoice", p.adapter.Adapt(p.handler.CreateInvoice))
		protected.Get("/orders/{orderID}/status", p.adapter.Adapt(p.handler.GetPaymentStatus))
	})
}
