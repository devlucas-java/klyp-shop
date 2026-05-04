package service

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dcart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type ShoppingCartService struct {
	log            *logger.Logger
	cartRepository repository.ShoppingCartRepository
	cartMapper     *mapper.ShoppingCartMapper
}

func NewShoppingCartService(
	log *logger.Logger,
	cartRepository repository.ShoppingCartRepository,
	cartMapper *mapper.ShoppingCartMapper,
) *ShoppingCartService {
	return &ShoppingCartService{
		log:            log,
		cartRepository: cartRepository,
		cartMapper:     cartMapper,
	}
}

func (s *ShoppingCartService) GetCart(auth *entity.User) (*dcart.ShoppingCartResponse, error) {
	cart, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to get shopping cart for user %s: %v", auth.ID, err)
		return nil, err
	}
	if cart == nil {
		cart = &entity.ShoppingCart{
			ID:        id.NewUUID(),
			UserID:    auth.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Items:     []*entity.ShoppingCartItem{},
		}
	}
	return s.cartMapper.ShoppingCartToResponse(cart), nil
}

func (s *ShoppingCartService) ClearCart(auth *entity.User) error {
	cart, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to get shopping cart for user %s: %v", auth.ID, err)
		return err
	}
	if cart == nil {
		return nil
	}
	return s.cartRepository.DeleteByID(cart.ID)
}

func updateCartTotals(cart *entity.ShoppingCart) {
	var total float64
	for _, item := range cart.Items {
		total += item.Subtotal()
	}
	cart.TotalBTC = total
	cart.UpdatedAt = time.Now()
}
