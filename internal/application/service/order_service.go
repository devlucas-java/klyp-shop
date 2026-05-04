package service

import (
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dorder"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type OrderService struct {
	log               *logger.Logger
	orderRepository   repository.OrderRepository
	userRepository    repository.UserRepository
	addressRepository repository.AddressRepository
	productRepository repository.ProductRepository
	orderMapper       *mapper.OrderMapper
}

func NewOrderService(
	log *logger.Logger,
	orderRepository repository.OrderRepository,
	userRepository repository.UserRepository,
	addressRepository repository.AddressRepository,
	productRepository repository.ProductRepository,
	orderMapper *mapper.OrderMapper,
) *OrderService {
	return &OrderService{
		log:               log,
		orderRepository:   orderRepository,
		userRepository:    userRepository,
		addressRepository: addressRepository,
		productRepository: productRepository,
		orderMapper:       orderMapper,
	}
}

func (s *OrderService) CreateOrder(auth *entity.User, req *dorder.CreateOrderRequest) (*dorder.OrderResponse, error) {
	s.log.Infof("Creating order for user %s", auth.ID)

	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user %s: %v", auth.ID, err)
		return nil, errors.ErrNotFound("User", err)
	}

	addressID, err := id.Parse(req.AddressID)
	if err != nil {
		return nil, errors.ErrInvalidUUID(err)
	}

	address, err := s.addressRepository.FindByID(addressID)
	if err != nil {
		s.log.Errorf("Failed to find address %s: %v", addressID, err)
		return nil, errors.ErrNotFound("Address", err)
	}

	if address.UserID != user.ID {
		s.log.Warnf("Address %s does not belong to user %s", addressID, user.ID)
		return nil, errors.ErrForbidden(fmt.Errorf("address does not belong to user"))
	}

	if len(req.Items) == 0 {
		return nil, errors.ErrBadRequest("at least one item is required", nil)
	}

	items := make([]entity.OrderItem, 0, len(req.Items))
	for _, itemReq := range req.Items {
		productID, err := id.Parse(itemReq.ProductID)
		if err != nil {
			return nil, errors.ErrInvalidUUID(err)
		}

		if itemReq.Quantity <= 0 {
			return nil, errors.ErrBadRequest("quantity must be greater than zero", nil)
		}

		product, err := s.productRepository.FindByID(productID)
		if err != nil {
			s.log.Errorf("Failed to find product %s: %v", productID, err)
			return nil, errors.ErrNotFound("Product", err)
		}

		items = append(items, *entity.NewOrderItem(productID, itemReq.Quantity, product.PriceBTC))
	}

	order := entity.NewOrder(user.ID, addressID, items)
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
	}

	createdOrder, err := s.orderRepository.Create(order)
	if err != nil {
		s.log.Errorf("Failed to create order for user %s: %v", auth.ID, err)
		return nil, err
	}

	s.log.Infof("Order created successfully for user %s", auth.ID)
	return s.orderMapper.OrderToResponse(createdOrder), nil
}

func (s *OrderService) GetOrder(auth *entity.User, orderID id.UUID) (*dorder.OrderResponse, error) {
	s.log.Infof("Getting order %s for user %s", orderID, auth.ID)

	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		s.log.Errorf("Failed to find order %s: %v", orderID, err)
		return nil, errors.ErrNotFound("Order", err)
	}

	if order.UserID != auth.ID {
		s.log.Warnf("Order %s does not belong to user %s", orderID, auth.ID)
		return nil, errors.ErrForbidden(fmt.Errorf("order does not belong to user"))
	}

	return s.orderMapper.OrderToResponse(order), nil
}

func (s *OrderService) ListUserOrders(auth *entity.User) ([]*dorder.OrderResponse, error) {
	s.log.Infof("Listing orders for user %s", auth.ID)

	orders, err := s.orderRepository.FindByUser(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to list orders for user %s: %v", auth.ID, err)
		return nil, err
	}

	return s.orderMapper.OrdersToResponses(orders), nil
}

func (s *OrderService) CancelOrder(auth *entity.User, orderID id.UUID) error {
	s.log.Infof("Cancelling order %s for user %s", orderID, auth.ID)

	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		s.log.Errorf("Failed to find order %s: %v", orderID, err)
		return errors.ErrNotFound("Order", err)
	}

	if order.UserID != auth.ID {
		s.log.Warnf("Order %s does not belong to user %s", orderID, auth.ID)
		return errors.ErrForbidden(fmt.Errorf("order does not belong to user"))
	}

	if order.Status != entity.OrderStatusPending {
		s.log.Warnf("Cannot cancel order %s in status %s", orderID, order.Status)
		return errors.ErrConflict("Order", fmt.Errorf("order cannot be cancelled in current status"))
	}

	order.Status = entity.OrderStatusCancelled
	_, err = s.orderRepository.Updates(order)
	if err != nil {
		s.log.Errorf("Failed to cancel order %s: %v", orderID, err)
		return err
	}

	s.log.Infof("Order %s cancelled successfully", orderID)
	return nil
}
