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
		// NotFound means no cart yet — return an empty one
		var domainErr *apperrors.DomainError
		if ok := isNotFound(err, &domainErr); ok {
			c = entity.NewShoppingCart(auth.ID)
			return s.cartMapper.ShoppingCartToResponse(c), nil
		}
		return nil, apperrors.Database(shoppingCartServiceTrace+".get_cart: failed to get shopping cart", err)
	}
	if c == nil {
		c = entity.NewShoppingCart(auth.ID)
	}
	return s.cartMapper.ShoppingCartToResponse(c), nil
}

func (s *ShoppingCartService) ClearCart(auth *entity.User) error {
	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		var domainErr *apperrors.DomainError
		if ok := isNotFound(err, &domainErr); ok {
			return nil
		}
		return apperrors.Database(shoppingCartServiceTrace+".clear_cart: failed to get shopping cart", err)
	}
	if c == nil {
		return nil
	}
	return s.cartRepository.DeleteByID(c.ID)
}

// isNotFound checks whether err is a DomainError with KindNotFound.
func isNotFound(err error, target **apperrors.DomainError) bool {
	var de *apperrors.DomainError
	if !errAs(err, &de) {
		return false
	}
	if target != nil {
		*target = de
	}
	return de.Kind == apperrors.KindNotFound
}

func errAs(err error, target **apperrors.DomainError) bool {
	if err == nil {
		return false
	}
	de, ok := err.(*apperrors.DomainError)
	if ok {
		*target = de
		return true
	}
	return false
}
