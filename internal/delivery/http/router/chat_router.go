package router

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/adapter"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/handler"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/socket"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type ChatRouter struct {
	jwtService     *jwt.JWTService
	chatHandler    *handler.ChatHandler
	wsHandler      *socket.ChatWSHandler
	log            *logger.Logger
	userRepository repository.UserRepository
}

func NewChatRouter(
	jwtService *jwt.JWTService,
	chatHandler *handler.ChatHandler,
	wsHandler *socket.ChatWSHandler,
	log *logger.Logger,
	userRepository repository.UserRepository,
) *ChatRouter {
	return &ChatRouter{
		jwtService:     jwtService,
		chatHandler:    chatHandler,
		wsHandler:      wsHandler,
		log:            log,
		userRepository: userRepository,
	}
}

func (c *ChatRouter) RegisterChatRoutes(r chi.Router) {
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware(c.jwtService, c.log, c.userRepository))
		protected.Post("/messages", adapter.Adapt(c.chatHandler.SendMessage))
		protected.Get("/messages/{peerID}", adapter.Adapt(c.chatHandler.GetConversation))
		protected.Handle("/ws", c.wsHandler)
	})
}
