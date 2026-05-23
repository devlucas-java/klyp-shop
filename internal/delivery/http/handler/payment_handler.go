package handler

import (
	"io"
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/utils"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
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
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	orderID, err := id.Parse(chi.URLParam(r, "orderID"))
	if err != nil {
		return apperrors.InvalidUUID(err)
	}
	res, err := h.paymentService.CreateInvoice(r.Context(), auth, orderID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusCreated, res)
	return nil
}

func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) error {
	auth, err := utils.GetAuth(r)
	if err != nil {
		return err
	}
	orderID, err := id.Parse(chi.URLParam(r, "orderID"))
	if err != nil {
		return apperrors.InvalidUUID(err)
	}
	res, err := h.paymentService.GetPaymentStatus(r.Context(), auth, orderID)
	if err != nil {
		return err
	}
	response.ResponseEntity(w, http.StatusOK, res)
	return nil
}

func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 65536))
	if err != nil {
		return apperrors.BadRequest("failed to read request body", err)
	}
	signature := r.Header.Get("BTCPay-Sig")
	if err := h.paymentService.HandleWebhook(r.Context(), body, signature); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
