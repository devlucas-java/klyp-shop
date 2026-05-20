package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
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
	orderPolicy       *policy.OrderPolicy
	metric            *metrics.Metric
}

func NewOrderService(
	log *logger.Logger,
	orderRepository repository.OrderRepository,
	userRepository repository.UserRepository,
	addressRepository repository.AddressRepository,
	productRepository repository.ProductRepository,
	orderMapper *mapper.OrderMapper,
	metric *metrics.Metric,
) *OrderService {
	return &OrderService{
		log:               log,
		orderRepository:   orderRepository,
		userRepository:    userRepository,
		addressRepository: addressRepository,
		productRepository: productRepository,
		orderMapper:       orderMapper,
		orderPolicy:       policy.NewOrderPolicy(),
		metric:            metric,
	}
}

func (s *OrderService) CreateOrder(auth *entity.User, req *order.CreateOrderRequest) (*order.OrderResponse, error) {
	s.log.Infof("Creating order for user %s", auth.ID)

	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	addressID, err := id.Parse(req.AddressID)
	if err != nil {
		return nil, errors.ErrInvalidUUID(err)
	}

	address, err := s.addressRepository.FindByID(addressID)
	if err != nil {
		return nil, errors.ErrNotFound("Address", err)
	}

	if err := s.orderPolicy.AddressBelongsToUser(address, user.ID); err != nil {
		return nil, err
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

		product, err := s.productRepository.FindByID(productID)
		if err != nil {
			return nil, errors.ErrNotFound("Product", err)
		}

		item, err := entity.NewOrderItem(productID, itemReq.Quantity, product.PriceBTC)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	order := entity.NewOrder(user.ID, addressID, items)
	order.SetOrderIDForItems()

	created, err := s.orderRepository.Create(order)
	if err != nil {
		return nil, errors.ErrDatabase("failed to create order", err)
	}

	s.metric.OrdersCreated.Inc()
	s.log.Infof("Order %s created for user %s", created.ID, auth.ID)
	return s.orderMapper.OrderToResponse(created), nil
}

func (s *OrderService) GetOrder(auth *entity.User, orderID id.UUID) (*order.OrderResponse, error) {
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return nil, errors.ErrNotFound("Order", err)
	}

	if err := s.orderPolicy.CanView(order, auth.ID); err != nil {
		return nil, err
	}

	return s.orderMapper.OrderToResponse(order), nil
}

func (s *OrderService) ListUserOrders(auth *entity.User) ([]*order.OrderResponse, error) {
	orders, err := s.orderRepository.FindByUser(auth.ID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to list orders", err)
	}

	return s.orderMapper.OrdersToResponses(orders), nil
}

func (s *OrderService) CancelOrder(auth *entity.User, orderID id.UUID) error {
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return errors.ErrNotFound("Order", err)
	}

	if err := s.orderPolicy.CanCancel(order, auth.ID); err != nil {
		return err
	}

	if err := order.CancelPending(); err != nil {
		return err
	}

	if _, err := s.orderRepository.Updates(order); err != nil {
		return errors.ErrDatabase("failed to cancel order", err)
	}

	s.metric.OrdersCancelled.Inc()
	s.log.Infof("Order %s cancelled by user %s", orderID, auth.ID)
	return nil
}
