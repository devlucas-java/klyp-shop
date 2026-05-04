package handler

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type ShoppingCartHandler struct {
	shoppingCartService *service.ShoppingCartService
	log                 *logger.Logger
}

func NewShoppingCartHandler(shoppingCartService *service.ShoppingCartService, log *logger.Logger) *ShoppingCartHandler {
	return &ShoppingCartHandler{shoppingCartService: shoppingCartService, log: log}
}

func (h *ShoppingCartHandler) GetCart(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	res, err := h.shoppingCartService.GetCart(auth)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *ShoppingCartHandler) ClearCart(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	if err := h.shoppingCartService.ClearCart(auth); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
