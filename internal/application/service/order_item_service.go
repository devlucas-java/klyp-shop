package service

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order_item"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const orderItemServiceTrace = "order_item_service.OrderItemService"

type OrderItemService struct {
	log                 *logger.Logger
	orderItemRepository repository.OrderItemRepository
	orderRepository     repository.OrderRepository
	productRepository   repository.ProductRepository
	orderMapper         *mapper.OrderMapper
}

func NewOrderItemService(
	log *logger.Logger,
	orderItemRepository repository.OrderItemRepository,
	orderRepository repository.OrderRepository,
	productRepository repository.ProductRepository,
	orderMapper *mapper.OrderMapper,
) *OrderItemService {
	return &OrderItemService{
		log:                 log,
		orderItemRepository: orderItemRepository,
		orderRepository:     orderRepository,
		productRepository:   productRepository,
		orderMapper:         orderMapper,
	}
}

func (s *OrderItemService) GetOrderItems(ctx context.Context, orderID id.UUID) ([]*dorderitem.OrderItemResponse, error) {
	_, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound(orderItemServiceTrace+".get_order_items: order not found", err)
	}

	items, err := s.orderItemRepository.FindByOrder(orderID)
	if err != nil {
		return nil, apperrors.Database(orderItemServiceTrace+".get_order_items: failed to find order items", err)
	}

	if len(items) == 0 {
		return []*dorderitem.OrderItemResponse{}, nil
	}

	responses := make([]*dorderitem.OrderItemResponse, len(items))
	for i, item := range items {
		responses[i] = &dorderitem.OrderItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			PriceBTC:  item.PriceBTC,
			Subtotal:  item.Subtotal(),
		}
	}

	return responses, nil
}

func (s *OrderItemService) GetOrderItem(ctx context.Context, orderID, itemID id.UUID) (*dorderitem.OrderItemResponse, error) {
	_, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.NotFound(orderItemServiceTrace+".get_order_item: order not found", err)
	}

	item, err := s.orderItemRepository.FindByID(itemID)
	if err != nil {
		return nil, apperrors.NotFound(orderItemServiceTrace+".get_order_item: order item not found", err)
	}

	if item.OrderID != orderID {
		return nil, apperrors.Forbidden(orderItemServiceTrace+".get_order_item: item does not belong to order", nil)
	}

	return &dorderitem.OrderItemResponse{
		ID:        item.ID.String(),
		ProductID: item.ProductID.String(),
		Quantity:  item.Quantity,
		PriceBTC:  item.PriceBTC,
		Subtotal:  item.Subtotal(),
	}, nil
}

func (s *OrderItemService) CalculateOrderTotal(orderID id.UUID) (float64, error) {
	items, err := s.orderItemRepository.FindByOrder(orderID)
	if err != nil {
		return 0, apperrors.Database(orderItemServiceTrace+".calculate_order_total: failed to find order items", err)
	}

	var total float64
	for _, item := range items {
		total += item.Subtotal()
	}

	return total, nil
}
