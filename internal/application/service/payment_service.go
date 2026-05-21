package service

import (
	"context"
	"math"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/payment"
	paymentDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/payment"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/client/port"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const paymentServiceTrace = "payment_service.PaymentService"

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

func (s *PaymentService) CreateInvoice(ctx context.Context, auth *entity.User, orderID id.UUID) (*payment.InvoiceResponse, error) {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound(paymentServiceTrace+".create_invoice: order not found", err)
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
		return nil, apperrors.Internal(paymentServiceTrace+".create_invoice: failed to create invoice", err)
	}

	pay := entity.NewBitcoinPayment(orderID, invoice.CheckoutLink, amountSats)
	pay.TxHash = invoice.ID

	saved, err := s.paymentRepository.Create(pay)
	if err != nil {
		return nil, apperrors.Database(paymentServiceTrace+".create_invoice: failed to save payment", err)
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

func (s *PaymentService) GetPaymentStatus(ctx context.Context, auth *entity.User, orderID id.UUID) (*payment.InvoiceResponse, error) {
	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound(paymentServiceTrace+".get_payment_status: order not found", err)
	}

	if err := order.EnsureOwnedBy(auth.ID); err != nil {
		return nil, err
	}

	pay, err := s.paymentRepository.FindByOrderID(orderID)
	if err != nil {
		return nil, apperrors.NotFound(paymentServiceTrace+".get_payment_status: payment not found", err)
	}

	if pay.TxHash != "" {
		invoice, err := s.paymentGateway.GetInvoice(pay.TxHash)
		if err == nil {
			s.syncPaymentStatus(pay, invoice.Status)
		}
	}

	return &paymentDTO.InvoiceResponse{
		PaymentID:     pay.ID.String(),
		OrderID:       pay.OrderID.String(),
		AmountSats:    pay.AmountSats,
		Status:        string(pay.Status),
		WalletAddress: pay.WalletAddress,
		InvoiceID:     pay.TxHash,
	}, nil
}

// HandleWebhook recebe o payload bruto e delega parse + validação de assinatura ao gateway.
func (s *PaymentService) HandleWebhook(ctx context.Context, rawBody []byte, signature string) error {
	event, err := s.paymentGateway.ParseWebhook(rawBody, signature)
	if err != nil {
		return apperrors.Unauthorized(paymentServiceTrace+".handle_webhook: invalid webhook signature", err)
	}

	s.log.Infof("PaymentService.HandleWebhook: type=%s invoiceID=%s orderID=%s",
		event.Type, event.InvoiceID, event.OrderID)

	switch event.Type {
	case "InvoiceSettled", "InvoicePaymentSettled":
		return s.handleInvoiceSettled(ctx, event)
	case "InvoiceExpired", "InvoiceInvalid":
		return s.handleInvoiceFailed(event)
	}

	return nil
}

func (s *PaymentService) handleInvoiceSettled(ctx context.Context, event *port.WebhookEvent) error {
	orderID, err := id.Parse(event.OrderID)
	if err != nil {
		return apperrors.InvalidUUID(paymentServiceTrace+".handle_invoice_settled: invalid order id", err)
	}

	pay, err := s.paymentRepository.FindByOrderID(orderID)
	if err != nil {
		return apperrors.NotFound(paymentServiceTrace+".handle_invoice_settled: payment not found", err)
	}

	pay.Confirm(event.InvoiceID)
	if _, err := s.paymentRepository.Save(pay); err != nil {
		return apperrors.Database(paymentServiceTrace+".handle_invoice_settled: failed to update payment", err)
	}

	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return apperrors.NotFound(paymentServiceTrace+".handle_invoice_settled: order not found", err)
	}

	order.MarkAsPaid()
	if _, err := s.orderRepository.Updates(ctx, order); err != nil {
		return apperrors.Database(paymentServiceTrace+".handle_invoice_settled: failed to update order", err)
	}

	s.metric.PaymentsSettled.Inc()
	s.log.Infof("PaymentService: order %s marked as paid", orderID)
	return nil
}

func (s *PaymentService) handleInvoiceFailed(event *port.WebhookEvent) error {
	orderID, err := id.Parse(event.OrderID)
	if err != nil {
		return apperrors.InvalidUUID(paymentServiceTrace+".handle_invoice_failed: invalid order id", err)
	}

	pay, err := s.paymentRepository.FindByOrderID(orderID)
	if err != nil {
		return nil
	}

	pay.Fail()
	if _, err := s.paymentRepository.Save(pay); err != nil {
		return apperrors.Database(paymentServiceTrace+".handle_invoice_failed: failed to update payment", err)
	}

	s.metric.PaymentsFailed.Inc()
	s.log.Infof("PaymentService: payment for order %s marked as failed", orderID)
	return nil
}

func (s *PaymentService) syncPaymentStatus(pay *entity.BitcoinPayment, gatewayStatus string) {
	switch gatewayStatus {
	case "Settled", "Complete":
		if !pay.IsConfirmed() {
			pay.Confirm(pay.TxHash)
			s.paymentRepository.Save(pay)
		}
	case "Expired", "Invalid":
		if pay.Status == entity.PaymentStatusPending {
			pay.Fail()
			s.paymentRepository.Save(pay)
		}
	}
}
