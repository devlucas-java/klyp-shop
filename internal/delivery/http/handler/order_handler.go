package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dorder"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type OrderHandler struct {
	orderService *service.OrderService
	log          *logger.Logger
}

func NewOrderHandler(orderService *service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, log: log}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	var req dorder.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest("invalid request payload", err)
	}
	res, err := h.orderService.CreateOrder(auth, &req)
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
		return errors.ErrInvalidUUID(err)
	}
	res, err := h.orderService.GetOrder(auth, orderID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)
	res, err := h.orderService.ListUserOrders(auth)
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
		return errors.ErrInvalidUUID(err)
	}
	if err := h.orderService.CancelOrder(auth, orderID); err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, nil)
	return nil
}
