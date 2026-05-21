package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/cart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const shoppingCartServiceTrace = "shopping_cart_service.ShoppingCartService"

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

func (s *ShoppingCartService) GetCart(auth *entity.User) (*cart.ShoppingCartResponse, error) {
	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		return nil, apperrors.Database(shoppingCartServiceTrace+".get_cart: failed to get shopping cart", err)
	}
	return s.cartMapper.ShoppingCartToResponse(c), nil
}

func (s *ShoppingCartService) ClearCart(auth *entity.User) error {
	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		return apperrors.Database(shoppingCartServiceTrace+".clear_cart: failed to get shopping cart", err)
	}
	return s.cartRepository.DeleteByID(c.ID)
}
