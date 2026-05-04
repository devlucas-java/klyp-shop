package service

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dcart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

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

func (s *ShoppingCartItemService) AddItem(auth *entity.User, req *dcart.AddShoppingCartItemRequest) (*dcart.ShoppingCartResponse, error) {
	if req.Quantity <= 0 {
		return nil, errors.ErrBadRequest("quantity must be greater than zero", nil)
	}

	productID, err := id.Parse(req.ProductID)
	if err != nil {
		return nil, errors.ErrInvalidUUID(err)
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		s.log.Errorf("Failed to find product %s: %v", productID, err)
		return nil, errors.ErrNotFound("Product", err)
	}

	cart, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to get shopping cart for user %s: %v", auth.ID, err)
		return nil, err
	}

	isNewCart := cart == nil
	if isNewCart {
		cart = &entity.ShoppingCart{
			ID:        id.NewUUID(),
			UserID:    auth.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Items:     []*entity.ShoppingCartItem{},
		}
	}

	cart.Items = append(cart.Items, entity.NewShoppingCartItem(cart.ID, product.ID, req.Quantity, product.PriceBTC))
	updateCartTotals(cart)

	if isNewCart {
		cart, err = s.cartRepository.Create(cart)
	} else {
		cart, err = s.cartRepository.Save(cart)
	}
	if err != nil {
		s.log.Errorf("Failed to save shopping cart for user %s: %v", auth.ID, err)
		return nil, err
	}

	return s.cartMapper.ShoppingCartToResponse(cart), nil
}

func (s *ShoppingCartItemService) UpdateItem(auth *entity.User, itemID id.UUID, req *dcart.UpdateShoppingCartItemRequest) (*dcart.ShoppingCartResponse, error) {
	if req.Quantity <= 0 {
		return nil, errors.ErrBadRequest("quantity must be greater than zero", nil)
	}

	cart, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to get shopping cart for user %s: %v", auth.ID, err)
		return nil, err
	}
	if cart == nil {
		return nil, errors.ErrNotFound("ShoppingCart", nil)
	}

	var itemToUpdate *entity.ShoppingCartItem
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemToUpdate = item
			break
		}
	}
	if itemToUpdate == nil {
		return nil, errors.ErrNotFound("ShoppingCartItem", nil)
	}

	itemToUpdate.Quantity = req.Quantity
	itemToUpdate.UpdatedAt = time.Now()
	cart.UpdatedAt = time.Now()
	updateCartTotals(cart)

	cart, err = s.cartRepository.Save(cart)
	if err != nil {
		s.log.Errorf("Failed to update shopping cart item %s: %v", itemID, err)
		return nil, err
	}

	return s.cartMapper.ShoppingCartToResponse(cart), nil
}

func (s *ShoppingCartItemService) RemoveItem(auth *entity.User, itemID id.UUID) error {
	cart, err := s.cartRepository.FindByUserID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to get shopping cart for user %s: %v", auth.ID, err)
		return err
	}
	if cart == nil {
		return errors.ErrNotFound("ShoppingCart", nil)
	}

	updatedItems := make([]*entity.ShoppingCartItem, 0, len(cart.Items))
	removed := false
	for _, item := range cart.Items {
		if item.ID == itemID {
			removed = true
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	if !removed {
		return errors.ErrNotFound("ShoppingCartItem", nil)
	}

	cart.Items = updatedItems
	updateCartTotals(cart)

	if len(cart.Items) == 0 {
		return s.cartRepository.DeleteByID(cart.ID)
	}

	_, err = s.cartRepository.Save(cart)
	if err != nil {
		s.log.Errorf("Failed to save shopping cart after removing item %s: %v", itemID, err)
		return err
	}

	return nil
}
