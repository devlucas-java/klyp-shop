package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/cart"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/utils"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

const shoppingCartItemHandlerTrace = "shopping_cart_item_handler.ShoppingCartItemHandler"

type ShoppingCartItemHandler struct {
	shoppingCartItemService *service.ShoppingCartItemService
	log                     *logger.Logger
}

func NewShoppingCartItemHandler(shoppingCartItemService *service.ShoppingCartItemService, log *logger.Logger) *ShoppingCartItemHandler {
	return &ShoppingCartItemHandler{shoppingCartItemService: shoppingCartItemService, log: log}
}

func (h *ShoppingCartItemHandler) AddItem(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}

	var req cart.AddShoppingCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest(shoppingCartItemHandlerTrace+".add_item: invalid request payload", err)
	}
	res, err := h.shoppingCartItemService.AddItem(auth, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *ShoppingCartItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}

	itemID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(shoppingCartItemHandlerTrace+".update_item: invalid item id", err)
	}
	var req cart.UpdateShoppingCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest(shoppingCartItemHandlerTrace+".update_item: invalid request payload", err)
	}
	res, err := h.shoppingCartItemService.UpdateItem(auth, itemID, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *ShoppingCartItemHandler) RemoveItem(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	itemID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(shoppingCartItemHandlerTrace+".remove_item: invalid item id", err)
	}
	if err := h.shoppingCartItemService.RemoveItem(auth, itemID); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
