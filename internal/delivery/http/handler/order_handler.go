package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/pkg/pagination"
	"github.com/go-chi/chi"
)

const orderHandlerTrace = "order_handler.OrderHandler"

type OrderHandler struct {
	orderService *service.OrderService
	log          *logger.Logger
}

func NewOrderHandler(orderService *service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, log: log}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.BadRequest(orderHandlerTrace+".create_order: invalid request payload", err)
	}
	res, err := h.orderService.CreateOrder(r.Context(), auth, &req)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	orderID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(orderHandlerTrace+".get_order_by_id: invalid order id", err)
	}
	res, err := h.orderService.GetOrder(r.Context(), auth, orderID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	p := pagination.ParsePagination(r)

	res, err := h.orderService.ListUserOrders(r.Context(), auth, p)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	orderID, err := id.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return apperrors.InvalidUUID(orderHandlerTrace+".cancel_order: invalid order id", err)
	}
	if err := h.orderService.CancelOrder(r.Context(), auth, orderID); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
