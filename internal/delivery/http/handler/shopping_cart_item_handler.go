package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dcart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type ShoppingCartItemHandler struct {
	shoppingCartItemService *service.ShoppingCartItemService
	log                     *logger.Logger
}

func NewShoppingCartItemHandler(shoppingCartItemService *service.ShoppingCartItemService, log *logger.Logger) *ShoppingCartItemHandler {
	return &ShoppingCartItemHandler{shoppingCartItemService: shoppingCartItemService, log: log}
}

func (h *ShoppingCartItemHandler) AddItem(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	var req dcart.AddShoppingCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrInvalidPayload(err)
	}
	res, err := h.shoppingCartItemService.AddItem(auth, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *ShoppingCartItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	itemID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	var req dcart.UpdateShoppingCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrInvalidPayload(err)
	}
	res, err := h.shoppingCartItemService.UpdateItem(auth, itemID, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *ShoppingCartItemHandler) RemoveItem(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	itemID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	if err := h.shoppingCartItemService.RemoveItem(auth, itemID); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
