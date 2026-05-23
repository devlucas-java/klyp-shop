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

type ShoppingCartItemService struct {
	log                  *logger.Logger
	cartRepository       repository.ShoppingCartRepository
	shoppingCartItemRepo repository.ShoppingCartItemRepository
	productRepo          repository.ProductRepository
	cartMapper           *mapper.ShoppingCartMapper
}

func NewShoppingCartItemService(
	log *logger.Logger,
	cartRepository repository.ShoppingCartRepository,
	shoppingCartItemRepo repository.ShoppingCartItemRepository,
	productRepo repository.ProductRepository,
	cartMapper *mapper.ShoppingCartMapper,
) *ShoppingCartItemService {
	return &ShoppingCartItemService{
		log:                  log,
		cartRepository:       cartRepository,
		shoppingCartItemRepo: shoppingCartItemRepo,
		productRepo:          productRepo,
		cartMapper:           cartMapper,
	}
}

func (s *ShoppingCartItemService) AddItem(auth *entity.User, req *cart.AddShoppingCartItemRequest) (*cart.ShoppingCartResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	productID, err := id.Parse(req.ProductID)
	if err != nil {
		return nil, apperrors.InvalidUUID(err)
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}

	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		return nil, err
	}

	existing := c.FindItemByProductID(product.ID)
	if existing != nil {
		if err := c.UpdateItemQuantity(existing.ID, existing.Quantity+req.Quantity); err != nil {
			return nil, err
		}
	} else {
		item, err := entity.NewShoppingCartItem(c.ID, product.ID, req.Quantity, product.PriceBTC)
		if err != nil {
			return nil, err
		}
		if err := c.AddItem(item); err != nil {
			return nil, err
		}
	}

	saved, err := s.cartRepository.Save(c)
	if err != nil {
		return nil, err
	}

	return s.cartMapper.ShoppingCartToResponse(saved), nil
}

func (s *ShoppingCartItemService) UpdateItem(auth *entity.User, itemID id.UUID, req *cart.UpdateShoppingCartItemRequest) (*cart.ShoppingCartResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		return nil, err
	}

	if err := c.UpdateItemQuantity(itemID, req.Quantity); err != nil {
		return nil, err
	}

	saved, err := s.cartRepository.Save(c)
	if err != nil {
		return nil, err
	}

	return s.cartMapper.ShoppingCartToResponse(saved), nil
}

func (s *ShoppingCartItemService) RemoveItem(auth *entity.User, itemID id.UUID) (*cart.ShoppingCartResponse, error) {
	c, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		return nil, err
	}

	if err := c.RemoveItem(itemID); err != nil {
		return nil, err
	}

	if err := s.shoppingCartItemRepo.DeleteByID(itemID); err != nil {
		return nil, err
	}

	saved, err := s.cartRepository.Save(c)
	if err != nil {
		return nil, err
	}

	return s.cartMapper.ShoppingCartToResponse(saved), nil
}
