package module

import (
	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/router"
	"github.com/devlucas-java/klyp-shop/internal/delivery/socket"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func InitChatModule(db *gorm.DB, log *logger.Logger, jwtService *jwt.JWTService, metric *metrics.Metric) (chi.Router, *socket.Hub) {
	chatRepository := database.NewChatDB(db)
	userRepository := database.NewUserDB(db)

	chatService := service.NewChatService(log, chatRepository, userRepository)
	hub := socket.NewHub(log, metric)

	chatHandler := handler.NewChatHandler(chatService, log)
	wsHandler := socket.NewChatWSHandler(hub, chatService, log)
	adapter := adapter.NewAdapter(log)
	chatRouter := router.NewChatRouter(jwtService, chatHandler, wsHandler, log, userRepository, adapter)

	r := chi.NewRouter()
	chatRouter.RegisterChatRoutes(r)
	return r, hub
}
