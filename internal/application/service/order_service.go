package service

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/pkg/pagination"
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

func (s *OrderService) CreateOrder(ctx context.Context, auth *entity.User, req *order.CreateOrderRequest) (*order.OrderResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, err
	}

	addressID, err := id.Parse(req.AddressID)
	if err != nil {
		return nil, apperrors.InvalidUUID(err)
	}

	address, err := s.addressRepository.FindByID(addressID)
	if err != nil {
		return nil, err
	}

	if err := s.orderPolicy.AddressBelongsToUser(address, user.ID); err != nil {
		return nil, err
	}

	items := make([]entity.OrderItem, 0, len(req.Items))
	for _, itemReq := range req.Items {
		productID, err := id.Parse(itemReq.ProductID)
		if err != nil {
			return nil, apperrors.InvalidUUID(err)
		}

		product, err := s.productRepository.FindByID(productID)
		if err != nil {
			return nil, err
		}

		item, err := entity.NewOrderItem(productID, itemReq.Quantity, product.PriceBTC)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	newOrder := entity.NewOrder(user.ID, addressID, items)
	newOrder.SetOrderIDForItems()

	created, err := s.orderRepository.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	s.metric.OrdersCreated.Inc()
	s.log.Infof("order %s created for user %s", created.ID, auth.ID)
	return s.orderMapper.OrderToResponse(created), nil
}

func (s *OrderService) GetOrder(ctx context.Context, auth *entity.User, orderID id.UUID) (*order.OrderResponse, error) {
	ord, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if err := s.orderPolicy.CanView(ord, auth.ID); err != nil {
		return nil, err
	}

	return s.orderMapper.OrderToResponse(ord), nil
}

func (s *OrderService) ListUserOrders(ctx context.Context, auth *entity.User, inputPagination pagination.InputPagination) (*order.OrdersPageResponse, error) {
	orders, total, err := s.orderRepository.FindByUserIDPaginated(ctx, auth.ID, inputPagination.Page, inputPagination.Size, inputPagination.Search)
	if err != nil {
		return nil, err
	}

	return &order.OrdersPageResponse{
		Pagination: buildPagination(inputPagination.Page, inputPagination.Size, total),
		Items:      s.orderMapper.OrdersToResponses(orders),
	}, nil
}

func buildPagination(page, size int, total int64) pagination.OutPutPagination {
	totalPages := int64((total + int64(size) - 1) / int64(size))
	if totalPages < 1 {
		totalPages = 1
	}

	return pagination.OutPutPagination{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}
}

func (s *OrderService) CancelOrder(ctx context.Context, auth *entity.User, orderID id.UUID) error {
	ord, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := s.orderPolicy.CanCancel(ord, auth.ID); err != nil {
		return err
	}

	if err := ord.CancelPending(); err != nil {
		return err
	}

	if _, err := s.orderRepository.Updates(ctx, ord); err != nil {
		return err
	}

	s.metric.OrdersCancelled.Inc()
	s.log.Infof("order %s cancelled by user %s", orderID, auth.ID)
	return nil
}
