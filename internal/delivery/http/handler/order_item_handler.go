package handler

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type OrderItemHandler struct {
	orderItemService *service.OrderItemService
	log              *logger.Logger
}

func NewOrderItemHandler(orderItemService *service.OrderItemService, log *logger.Logger) *OrderItemHandler {
	return &OrderItemHandler{orderItemService: orderItemService, log: log}
}

func (h *OrderItemHandler) GetOrderItems(w http.ResponseWriter, r *http.Request) error {
	orderID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	res, err := h.orderItemService.GetOrderItems(orderID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *OrderItemHandler) GetOrderItem(w http.ResponseWriter, r *http.Request) error {
	orderID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	itemID, err := id.Parse(chi.URLParam(r, "itemId"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}
	res, err := h.orderItemService.GetOrderItem(orderID, itemID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
