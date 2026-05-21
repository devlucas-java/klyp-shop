package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/cart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const shoppingCartItemServiceTrace = "shopping_cart_item_service.ShoppingCartItemService"

type ShoppingCartItemService struct {
	log            *logger.Logger
	cartRepository repository.ShoppingCartRepository
	productRepo    repository.ProductRepository
	cartMapper     *mapper.ShoppingCartMapper
}

func NewShoppingCartItemService(
	log *logger.Logger,
	cartRepository repository.ShoppingCartRepository,
	productRepo repository.ProductRepository,
	cartMapper *mapper.ShoppingCartMapper,
) *ShoppingCartItemService {
	return &ShoppingCartItemService{
		log:            log,
		cartRepository: cartRepository,
		productRepo:    productRepo,
		cartMapper:     cartMapper,
	}
}

func (s *ShoppingCartItemService) AddItem(auth *entity.User, req *cart.AddShoppingCartItemRequest) (*cart.ShoppingCartResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperrors.BadRequest(shoppingCartItemServiceTrace+".add_item: quantity must be greater than zero", nil)
	}

	productID, err := id.Parse(req.ProductID)
	if err != nil {
		return nil, apperrors.InvalidUUID(shoppingCartItemServiceTrace+".add_item: invalid product id", err)
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, apperrors.NotFound(shoppingCartItemServiceTrace+".add_item: product not found", err)
	}

	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		var de *apperrors.DomainError
		if !errAs(err, &de) || de.Kind != apperrors.KindNotFound {
			return nil, apperrors.Database(shoppingCartItemServiceTrace+".add_item: failed to get shopping cart", err)
		}
		c = nil
	}

	isNewCart := c == nil
	if isNewCart {
		c = entity.NewShoppingCart(auth.ID)
	}

	item, err := entity.NewShoppingCartItem(c.ID, product.ID, req.Quantity, product.PriceBTC)
	if err != nil {
		return nil, err
	}

	if err := c.AddItem(item); err != nil {
		return nil, err
	}

	if isNewCart {
		c, err = s.cartRepository.Create(c)
	} else {
		c, err = s.cartRepository.Save(c)
	}
	if err != nil {
		return nil, apperrors.Database(shoppingCartItemServiceTrace+".add_item: failed to save shopping cart", err)
	}

	return s.cartMapper.ShoppingCartToResponse(c), nil
}

func (s *ShoppingCartItemService) UpdateItem(auth *entity.User, itemID id.UUID, req *cart.UpdateShoppingCartItemRequest) (*cart.ShoppingCartResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		return nil, apperrors.Database(shoppingCartItemServiceTrace+".update_item: failed to get shopping cart", err)
	}
	if c == nil {
		return nil, apperrors.NotFound(shoppingCartItemServiceTrace+".update_item: shopping cart not found", nil)
	}

	if err := c.UpdateItemQuantity(itemID, req.Quantity); err != nil {
		return nil, err
	}

	c, err = s.cartRepository.Save(c)
	if err != nil {
		return nil, apperrors.Database(shoppingCartItemServiceTrace+".update_item: failed to update shopping cart", err)
	}

	return s.cartMapper.ShoppingCartToResponse(c), nil
}

func (s *ShoppingCartItemService) RemoveItem(auth *entity.User, itemID id.UUID) error {
	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		var de *apperrors.DomainError
		if !errAs(err, &de) || de.Kind != apperrors.KindNotFound {
			return apperrors.Database(shoppingCartItemServiceTrace+".remove_item: failed to get shopping cart", err)
		}
		return nil
	}
	if c == nil {
		return nil
	}

	if err := c.RemoveItem(itemID); err != nil {
		return err
	}

	if len(c.Items) == 0 {
		return s.cartRepository.DeleteByID(c.ID)
	}

	_, err = s.cartRepository.Save(c)
	if err != nil {
		return apperrors.Database(shoppingCartItemServiceTrace+".remove_item: failed to save shopping cart", err)
	}

	return nil
}
