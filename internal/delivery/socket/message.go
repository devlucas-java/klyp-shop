package socket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dchat"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// messageProcessor handles reading, validating, persisting and routing messages.
type messageProcessor struct {
	hub         *Hub
	chatService *service.ChatService
	log         *logger.Logger
}

func newMessageProcessor(hub *Hub, chatService *service.ChatService, log *logger.Logger) *messageProcessor {
	return &messageProcessor{hub: hub, chatService: chatService, log: log}
}

func (p *messageProcessor) readPump(ctx context.Context, conn *websocket.Conn, auth *entity.User, cancel context.CancelFunc) {
	defer cancel()

	for {
		readCtx, readCancel := context.WithTimeout(ctx, readTimeout)
		var req dchat.SendMessageRequest
		err := wsjson.Read(readCtx, conn, &req)
		readCancel()

		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				p.log.Infof("message.go: user %s disconnected normally", auth.ID)
			} else {
				p.log.Errorf("message.go: read error user=%s: %v", auth.ID, err)
			}
			return
		}

		if err := req.Validate(); err != nil {
			sendError(ctx, conn, err.Error(), p.log)
			continue
		}

		msg, err := p.chatService.SendMessage(auth, &req)
		if err != nil {
			sendError(ctx, conn, err.Error(), p.log)
			continue
		}

		wsMsg := dchat.WSMessage{Type: "message", Payload: *msg}
		data, err := json.Marshal(wsMsg)
		if err != nil {
			p.log.Errorf("message.go: marshal error: %v", err)
			continue
		}

		// Echo back to sender
		writeCtx, writeCancel := context.WithTimeout(ctx, writeTimeout)
		conn.Write(writeCtx, websocket.MessageText, data)
		writeCancel()

		// Deliver to receiver if online
		receiverID, err := id.Parse(req.ReceiverID)
		if err == nil && p.hub.IsOnline(receiverID) {
			p.hub.Send(receiverID, data)
		}
	}
}

// writePump drains the client's send channel and writes to the WebSocket.
func writePump(ctx context.Context, conn *websocket.Conn, client *Client, log *logger.Logger) {
	for {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
			if err := conn.Write(writeCtx, websocket.MessageText, msg); err != nil {
				log.Errorf("message.go: write error user=%s: %v", client.UserID, err)
				cancel()
				return
			}
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func sendError(ctx context.Context, conn *websocket.Conn, msg string, log *logger.Logger) {
	data, _ := json.Marshal(map[string]string{"type": "error", "message": msg})
	writeCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := conn.Write(writeCtx, websocket.MessageText, data); err != nil {
		log.Errorf("message.go: sendError write failed: %v", err)
	}
}
