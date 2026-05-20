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

type DashboardRouter struct {
	jwtService       *jwt.JWTService
	dashboardHandler *handler.DashboardHandler
	log              *logger.Logger
	userRepository   repository.UserRepository
}

func NewDashboardRouter(
	jwtService *jwt.JWTService,
	h *handler.DashboardHandler,
	log *logger.Logger,
	userRepository repository.UserRepository,
) *DashboardRouter {
	return &DashboardRouter{
		jwtService:       jwtService,
		dashboardHandler: h,
		log:              log,
		userRepository:   userRepository,
	}
}

func (d *DashboardRouter) RegisterDashboardRoutes(r chi.Router) {
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(d.jwtService, d.log, d.userRepository))

		protected.Group(func(sellerOnly chi.Router) {
			sellerOnly.Use(middleware.RoleMiddleware([]enums.Role{enums.SELLER, enums.ADMIN}))
			sellerOnly.Get("/seller", adapter.Adapt(d.dashboardHandler.GetSellerDashboard))
		})

		protected.Group(func(adminOnly chi.Router) {
			adminOnly.Use(middleware.RoleMiddleware([]enums.Role{enums.ADMIN}))
			adminOnly.Get("/admin", adapter.Adapt(d.dashboardHandler.GetAdminDashboard))
		})
	})
}
