package handler

import (
	"io"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/middleware"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/go-chi/chi"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	log            *logger.Logger
}

func NewPaymentHandler(paymentService *service.PaymentService, log *logger.Logger) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, log: log}
}

func (h *PaymentHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	orderID, err := id.Parse(chi.URLParam(r, "orderID"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	res, err := h.paymentService.CreateInvoice(auth, orderID)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) error {
	auth := r.Context().Value(middleware.AuthKey).(*entity.User)

	orderID, err := id.Parse(chi.URLParam(r, "orderID"))
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	res, err := h.paymentService.GetPaymentStatus(auth, orderID)
	if err != nil {
		return err
	}

	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 65536))
	if err != nil {
		return errors.ErrBadRequest("failed to read body", err)
	}

	signature := r.Header.Get("BTCPay-Sig")

	if err := h.paymentService.HandleWebhook(body, signature); err != nil {
		h.log.Errorf("PaymentHandler.Webhook: %v", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
