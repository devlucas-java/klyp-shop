package socket

import (
	"context"
	"net/http"
	"time"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/coder/websocket"
)

const (
	writeTimeout = 10 * time.Second
	readTimeout  = 60 * time.Second
	maxMsgBytes  = 8192
)

// ChatWSHandler manages WebSocket lifecycle: accept, register, pump, cleanup.
type ChatWSHandler struct {
	hub         *Hub
	chatService *service.ChatService
	log         *logger.Logger
}

func NewChatWSHandler(hub *Hub, chatService *service.ChatService, log *logger.Logger) *ChatWSHandler {
	return &ChatWSHandler{hub: hub, chatService: chatService, log: log}
}

func (h *ChatWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth, ok := r.Context().Value(middleware.AuthKey).(*entity.User)
	if !ok || auth == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		h.log.Errorf("chat.go: accept error user=%s: %v", auth.ID, err)
		return
	}
	conn.SetReadLimit(maxMsgBytes)

	client := &Client{
		UserID: auth.ID,
		Send:   make(chan []byte, 64),
	}
	h.hub.Register(client)
	defer h.hub.Unregister(auth.ID)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	processor := newMessageProcessor(h.hub, h.chatService, h.log)

	go writePump(ctx, conn, client, h.log)
	processor.readPump(ctx, conn, auth, cancel)
}
