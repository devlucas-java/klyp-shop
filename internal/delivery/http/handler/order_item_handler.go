package handler

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

const orderItemHandlerTrace = "order_item_handler.OrderItemHandler"

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
		return apperrors.InvalidUUID(orderItemHandlerTrace+".get_order_items: invalid order id", err)
	}
	res, err := h.orderItemService.GetOrderItems(r.Context(), orderID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *OrderItemHandler) GetOrderItem(w http.ResponseWriter, r *http.Request) error {
	orderID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(orderItemHandlerTrace+".get_order_item: invalid order id", err)
	}
	itemID, err := id.Parse(chi.URLParam(r, "itemId"))
	if err != nil {
		return apperrors.InvalidUUID(orderItemHandlerTrace+".get_order_item: invalid item id", err)
	}
	res, err := h.orderItemService.GetOrderItem(r.Context(), orderID, itemID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}
