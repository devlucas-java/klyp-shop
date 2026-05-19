package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dpayment"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/btcpay"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type PaymentService struct {
	log               *logger.Logger
	paymentRepository repository.BitcoinPaymentRepository
	orderRepository   repository.OrderRepository
	btcpayClient      *btcpay.Client
	webhookSecret     string
}

func NewPaymentService(
	log *logger.Logger,
	paymentRepository repository.BitcoinPaymentRepository,
	orderRepository repository.OrderRepository,
	btcpayClient *btcpay.Client,
	webhookSecret string,
) *PaymentService {
	return &PaymentService{
		log:               log,
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
		btcpayClient:      btcpayClient,
		webhookSecret:     webhookSecret,
	}
}

func (s *PaymentService) CreateInvoice(auth *entity.User, orderID id.UUID) (*dpayment.InvoiceResponse, error) {
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return nil, errors.ErrNotFound("Order", err)
	}

	if err := order.CanBePaidBy(auth.ID); err != nil {
		return nil, err
	}

	existing, _ := s.paymentRepository.FindByOrderID(orderID)
	if existing != nil {
		return &dpayment.InvoiceResponse{
			PaymentID:     existing.ID.String(),
			OrderID:       existing.OrderID.String(),
			AmountBTC:     existing.AmountBTC,
			Status:        string(existing.Status),
			WalletAddress: existing.WalletAddress,
		}, nil
	}

	invoice, err := s.btcpayClient.CreateInvoice(orderID.String(), order.TotalBTC)
	if err != nil {
		s.log.Errorf("PaymentService.CreateInvoice btcpay error: %v", err)
		return nil, errors.ErrInternal("failed to create btcpay invoice", err)
	}

	payment := entity.NewBitcoinPayment(orderID, invoice.CheckoutLink, order.TotalBTC)
	payment.TxHash = invoice.ID

	saved, err := s.paymentRepository.Create(payment)
	if err != nil {
		return nil, errors.ErrDatabase("failed to save payment", err)
	}

	return &dpayment.InvoiceResponse{
		PaymentID:     saved.ID.String(),
		OrderID:       saved.OrderID.String(),
		AmountBTC:     saved.AmountBTC,
		Status:        string(saved.Status),
		WalletAddress: saved.WalletAddress,
		CheckoutURL:   invoice.CheckoutLink,
		InvoiceID:     invoice.ID,
	}, nil
}

func (s *PaymentService) GetPaymentStatus(auth *entity.User, orderID id.UUID) (*dpayment.InvoiceResponse, error) {
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return nil, errors.ErrNotFound("Order", err)
	}

	if err := order.EnsureOwnedBy(auth.ID); err != nil {
		return nil, err
	}

	payment, err := s.paymentRepository.FindByOrderID(orderID)
	if err != nil {
		return nil, errors.ErrNotFound("Payment", err)
	}

	if payment.TxHash != "" {
		invoice, err := s.btcpayClient.GetInvoice(payment.TxHash)
		if err == nil {
			s.syncPaymentStatus(payment, invoice.Status)
		}
	}

	return &dpayment.InvoiceResponse{
		PaymentID:     payment.ID.String(),
		OrderID:       payment.OrderID.String(),
		AmountBTC:     payment.AmountBTC,
		Status:        string(payment.Status),
		WalletAddress: payment.WalletAddress,
		InvoiceID:     payment.TxHash,
	}, nil
}

func (s *PaymentService) HandleWebhook(rawBody []byte, signature string) error {
	if !s.verifyWebhookSignature(rawBody, signature) {
		s.log.Warnf("PaymentService.HandleWebhook: invalid signature")
		return errors.ErrUnauthorized(fmt.Errorf("invalid webhook signature"))
	}

	var payload btcpay.WebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return errors.ErrBadRequest("invalid webhook payload", err)
	}

	s.log.Infof("PaymentService.HandleWebhook: type=%s invoiceID=%s orderID=%s",
		payload.Type, payload.InvoiceID, payload.Metadata.OrderID)

	switch payload.Type {
	case "InvoiceSettled", "InvoicePaymentSettled":
		return s.handleInvoiceSettled(payload)
	case "InvoiceExpired", "InvoiceInvalid":
		return s.handleInvoiceFailed(payload)
	}

	return nil
}

func (s *PaymentService) handleInvoiceSettled(payload btcpay.WebhookPayload) error {
	orderID, err := id.Parse(payload.Metadata.OrderID)
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	payment, err := s.paymentRepository.FindByOrderID(orderID)
	if err != nil {
		return errors.ErrNotFound("Payment", err)
	}

	payment.Confirm(payload.InvoiceID)
	if _, err := s.paymentRepository.Save(payment); err != nil {
		return errors.ErrDatabase("failed to update payment", err)
	}

	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return errors.ErrNotFound("Order", err)
	}

	order.MarkAsPaid()
	if _, err := s.orderRepository.Updates(order); err != nil {
		return errors.ErrDatabase("failed to update order", err)
	}

	s.log.Infof("PaymentService: order %s marked as paid", orderID)
	return nil
}

func (s *PaymentService) handleInvoiceFailed(payload btcpay.WebhookPayload) error {
	orderID, err := id.Parse(payload.Metadata.OrderID)
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	payment, err := s.paymentRepository.FindByOrderID(orderID)
	if err != nil {
		return nil
	}

	payment.Fail()
	if _, err := s.paymentRepository.Save(payment); err != nil {
		return errors.ErrDatabase("failed to update payment", err)
	}

	s.log.Infof("PaymentService: payment for order %s marked as failed", orderID)
	return nil
}

func (s *PaymentService) syncPaymentStatus(payment *entity.BitcoinPayment, btcpayStatus string) {
	switch btcpayStatus {
	case "Settled", "Complete":
		if !payment.IsConfirmed() {
			payment.Confirm(payment.TxHash)
			s.paymentRepository.Save(payment)
		}
	case "Expired", "Invalid":
		if payment.Status == entity.PaymentStatusPending {
			payment.Fail()
			s.paymentRepository.Save(payment)
		}
	}
}

func (s *PaymentService) verifyWebhookSignature(body []byte, signature string) bool {
	if s.webhookSecret == "" {
		return true
	}
	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
