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
	adapter        *adapter.Adapter
}

func NewChatRouter(
	jwt *jwt.JWTService,
	ch *handler.ChatHandler,
	ws *socket.ChatWSHandler,
	l *logger.Logger,
	ur repository.UserRepository,
	a *adapter.Adapter,
) *ChatRouter {
	return &ChatRouter{
		jwtService:     jwt,
		chatHandler:    ch,
		wsHandler:      ws,
		log:            l,
		userRepository: ur,
		adapter:        a,
	}
}

func (c *ChatRouter) RegisterChatRoutes(r chi.Router) {
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JwtMiddleware(c.jwtService, c.log, c.userRepository))
		protected.Post("/messages", c.adapter.Adapt(c.chatHandler.SendMessage))
		protected.Get("/messages/{peerID}", c.adapter.Adapt(c.chatHandler.GetConversation))
		protected.Handle("/ws", c.wsHandler)
	})
}
