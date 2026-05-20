package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/chat"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type ChatHandler struct {
	chatService *service.ChatService
	log         *logger.Logger
}

func NewChatHandler(chatService *service.ChatService, log *logger.Logger) *ChatHandler {
	return &ChatHandler{chatService: chatService, log: log}
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	var req chat.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrInvalidPayload(err)
	}
	if err := req.Validate(); err != nil {
		return errors.ErrBadRequest(err.Error(), nil)
	}

	res, err := h.chatService.SendMessage(auth, &req)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *ChatHandler) GetConversation(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	peerID, err := id.Parse(chi.URLParam(r, "peerID"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	msgs, err := h.chatService.GetConversation(auth, peerID, limit, offset)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusOK, msgs)
	return nil
}
