package service

import (
	dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dorder_item"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
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

func (s *OrderItemService) GetOrderItems(orderID id.UUID) ([]*dorderitem.OrderItemResponse, error) {
	s.log.Infof("Getting items for order %s", orderID)

	_, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		s.log.Errorf("Failed to find order %s: %v", orderID, err)
		return nil, errors.ErrNotFound("Order", err)
	}

	items, err := s.orderItemRepository.FindByOrder(orderID)
	if err != nil {
		s.log.Errorf("Failed to find items for order %s: %v", orderID, err)
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

func (s *OrderItemService) GetOrderItem(orderID, itemID id.UUID) (*dorderitem.OrderItemResponse, error) {
	s.log.Infof("Getting item %s from order %s", itemID, orderID)

	_, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		s.log.Errorf("Failed to find order %s: %v", orderID, err)
		return nil, errors.ErrNotFound("Order", err)
	}

	item, err := s.orderItemRepository.FindByID(itemID)
	if err != nil {
		s.log.Errorf("Failed to find item %s: %v", itemID, err)
		return nil, errors.ErrNotFound("OrderItem", err)
	}

	if item.OrderID != orderID {
		s.log.Warnf("Item %s does not belong to order %s", itemID, orderID)
		return nil, errors.ErrForbidden(nil)
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
	s.log.Infof("Calculating total for order %s", orderID)

	items, err := s.orderItemRepository.FindByOrder(orderID)
	if err != nil {
		s.log.Errorf("Failed to find items for order %s: %v", orderID, err)
		return 0, err
	}

	var total float64
	for _, item := range items {
		total += item.Subtotal()
	}

	return total, nil
}
