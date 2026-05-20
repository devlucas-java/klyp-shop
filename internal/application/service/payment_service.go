package service

import (
	"fmt"
	"math"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/payment"
	paymentDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/payment"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/client/port"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const satsPerBTC = int64(100_000_000)

// btcToSats converte um valor em BTC (float64) para satoshis (int64).
func btcToSats(btc float64) int64 {
	return int64(math.Round(btc * float64(satsPerBTC)))
}

type PaymentService struct {
	log               *logger.Logger
	paymentRepository repository.BitcoinPaymentRepository
	orderRepository   repository.OrderRepository
	paymentGateway    port.PaymentGateway
	metric            *metrics.Metric
}

func NewPaymentService(
	log *logger.Logger,
	paymentRepository repository.BitcoinPaymentRepository,
	orderRepository repository.OrderRepository,
	paymentGateway port.PaymentGateway,
	metric *metrics.Metric,
) *PaymentService {
	return &PaymentService{
		log:               log,
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
		paymentGateway:    paymentGateway,
		metric:            metric,
	}
}

func (s *PaymentService) CreateInvoice(auth *entity.User, orderID id.UUID) (*payment.InvoiceResponse, error) {
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return nil, errors.ErrNotFound("Order", err)
	}

	if err := order.CanBePaidBy(auth.ID); err != nil {
		return nil, err
	}

	// Retorna invoice existente sem criar duplicata
	existing, _ := s.paymentRepository.FindByOrderID(orderID)
	if existing != nil {
		return &payment.InvoiceResponse{
			PaymentID:     existing.ID.String(),
			OrderID:       existing.OrderID.String(),
			AmountSats:    existing.AmountSats,
			Status:        string(existing.Status),
			WalletAddress: existing.WalletAddress,
		}, nil
	}

	amountSats := btcToSats(order.TotalBTC)

	invoice, err := s.paymentGateway.CreateInvoice(orderID.String(), amountSats)
	if err != nil {
		s.log.Errorf("PaymentService.CreateInvoice gateway error: %v", err)
		return nil, errors.ErrInternal("failed to create invoice", err)
	}

	payment := entity.NewBitcoinPayment(orderID, invoice.CheckoutLink, amountSats)
	payment.TxHash = invoice.ID

	saved, err := s.paymentRepository.Create(payment)
	if err != nil {
		return nil, errors.ErrDatabase("failed to save payment", err)
	}

	s.metric.PaymentsCreated.Inc()

	return &paymentDTO.InvoiceResponse{
		PaymentID:     saved.ID.String(),
		OrderID:       saved.OrderID.String(),
		AmountSats:    saved.AmountSats,
		Status:        string(saved.Status),
		WalletAddress: saved.WalletAddress,
		CheckoutURL:   invoice.CheckoutLink,
		InvoiceID:     invoice.ID,
	}, nil
}

func (s *PaymentService) GetPaymentStatus(auth *entity.User, orderID id.UUID) (*payment.InvoiceResponse, error) {
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
		invoice, err := s.paymentGateway.GetInvoice(payment.TxHash)
		if err == nil {
			s.syncPaymentStatus(payment, invoice.Status)
		}
	}

	return &paymentDTO.InvoiceResponse{
		PaymentID:     payment.ID.String(),
		OrderID:       payment.OrderID.String(),
		AmountSats:    payment.AmountSats,
		Status:        string(payment.Status),
		WalletAddress: payment.WalletAddress,
		InvoiceID:     payment.TxHash,
	}, nil
}

// HandleWebhook recebe o payload bruto e delega parse + validação de assinatura ao gateway.
func (s *PaymentService) HandleWebhook(rawBody []byte, signature string) error {
	event, err := s.paymentGateway.ParseWebhook(rawBody, signature)
	if err != nil {
		s.log.Warnf("PaymentService.HandleWebhook: %v", err)
		return errors.ErrUnauthorized(fmt.Errorf("invalid webhook: %w", err))
	}

	s.log.Infof("PaymentService.HandleWebhook: type=%s invoiceID=%s orderID=%s",
		event.Type, event.InvoiceID, event.OrderID)

	switch event.Type {
	case "InvoiceSettled", "InvoicePaymentSettled":
		return s.handleInvoiceSettled(event)
	case "InvoiceExpired", "InvoiceInvalid":
		return s.handleInvoiceFailed(event)
	}

	return nil
}

func (s *PaymentService) handleInvoiceSettled(event *port.WebhookEvent) error {
	orderID, err := id.Parse(event.OrderID)
	if err != nil {
		return errors.ErrInvalidUUID(err)
	}

	payment, err := s.paymentRepository.FindByOrderID(orderID)
	if err != nil {
		return errors.ErrNotFound("Payment", err)
	}

	payment.Confirm(event.InvoiceID)
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

	s.metric.PaymentsSettled.Inc()
	s.log.Infof("PaymentService: order %s marked as paid", orderID)
	return nil
}

func (s *PaymentService) handleInvoiceFailed(event *port.WebhookEvent) error {
	orderID, err := id.Parse(event.OrderID)
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

	s.metric.PaymentsFailed.Inc()
	s.log.Infof("PaymentService: payment for order %s marked as failed", orderID)
	return nil
}

func (s *PaymentService) syncPaymentStatus(payment *entity.BitcoinPayment, gatewayStatus string) {
	switch gatewayStatus {
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
