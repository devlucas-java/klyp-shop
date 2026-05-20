package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitDashboardModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService) chi.Router {
	userRepository := database.NewUserDB(db, log)
	orderRepository := database.NewOrderDB(db, log)
	dashboardRepository := database.NewDashboardDB(db, log)

	dashboardService := service.NewDashboardService(log, userRepository, orderRepository, dashboardRepository)
	dashboardHandler := handler.NewDashboardHandler(dashboardService, log)
	dashboardRouter := router.NewDashboardRouter(jwtService, dashboardHandler, log, userRepository)

	r := chi.NewRouter()
	dashboardRouter.RegisterDashboardRoutes(r)
	return r
}
