package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/cart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
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

func (s *ShoppingCartService) GetCart(auth *entity.User) (*cart.ShoppingCartResponse, error) {
	cart, err := s.cartRepository.FindByUserID(auth.ID)

	if err != nil {
		return nil, err
	}

	return s.cartMapper.ShoppingCartToResponse(cart), nil
}

func (s *ShoppingCartService) ClearCart(auth *entity.User) error {
	cart, err := s.cartRepository.FindByUserID(auth.ID)

	if err != nil {
		return err
	}

	if err := s.cartRepository.DeleteByID(cart.ID); err != nil {
		return err
	}

	_, err = s.cartRepository.Create(entity.NewShoppingCart(auth.ID))
	return err
}
