package service

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order_item"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type OrderItemService struct {
	log                 *logger.Logger
	orderItemRepository repository.OrderItemRepository
	orderRepository     repository.OrderRepository
	productRepository   repository.ProductRepository
	orderMapper         *mapper.OrderMapper
	orderPolicy         *policy.OrderPolicy
}

func NewOrderItemService(
	log *logger.Logger,
	orderItemRepository repository.OrderItemRepository,
	orderRepository repository.OrderRepository,
	productRepository repository.ProductRepository,
	orderMapper *mapper.OrderMapper,
	orderPolicy *policy.OrderPolicy,
) *OrderItemService {
	return &OrderItemService{
		log:                 log,
		orderItemRepository: orderItemRepository,
		orderRepository:     orderRepository,
		productRepository:   productRepository,
		orderMapper:         orderMapper,
		orderPolicy:         orderPolicy,
	}
}

func (s *OrderItemService) GetOrderItems(ctx context.Context, orderID id.UUID) ([]*dorderitem.OrderItemResponse, error) {
	_, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items, err := s.orderItemRepository.FindByOrder(orderID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	item, err := s.orderItemRepository.FindByID(itemID)
	if err != nil {
		return nil, err
	}

	err = s.orderPolicy.ItemBelongsToOrder(item, orderID)
	if err != nil {
		return nil, err
	}

	return &dorderitem.OrderItemResponse{
		ID:        item.ID.String(),
		ProductID: item.ProductID.String(),
		Quantity:  item.Quantity,
		PriceBTC:  item.PriceBTC,
		Subtotal:  item.Subtotal(),
	}, nil
}

func (s *OrderItemService) CalculateOrderTotal(orderID id.UUID) (int64, error) {
	items, err := s.orderItemRepository.FindByOrder(orderID)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, item := range items {
		total += item.Subtotal()
	}

	return total, nil
}
