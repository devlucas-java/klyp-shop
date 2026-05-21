package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const shoppingCartEntity = "shopping_cart_entity.ShoppingCart"

type ShoppingCart struct {
	ID        id.UUID             `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time           `gorm:"autoCreateTime"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime"`
	UserID    id.UUID             `gorm:"index;not null"`
	Items     []*ShoppingCartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;"`
	TotalBTC  float64             `gorm:"type:decimal(18,8);not null"`
}

func NewShoppingCart(userID id.UUID) *ShoppingCart {
	now := time.Now()
	return &ShoppingCart{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		Items:     []*ShoppingCartItem{},
	}
}

func (c *ShoppingCart) AddItem(item *ShoppingCartItem) error {
	if item == nil {
		return apperrors.BadRequest(shoppingCartEntity+".add_item: invalid shopping cart item", nil)
	}
	if item.CartID != c.ID {
		return apperrors.BadRequest(shoppingCartEntity+".add_item: shopping cart item cart mismatch", nil)
	}

	c.Items = append(c.Items, item)
	c.RecalculateTotal()
	return nil
}

func (c *ShoppingCart) UpdateItemQuantity(itemID id.UUID, quantity int) error {
	if quantity <= 0 {
		return apperrors.BadRequest(shoppingCartEntity+".update_item_quantity: quantity must be greater than zero", nil)
	}

	item := c.FindItem(itemID)
	if item == nil {
		return apperrors.NotFound(shoppingCartEntity+".update_item_quantity: shopping cart item not found", nil)
	}

	if err := item.SetQuantity(quantity); err != nil {
		return err
	}
	c.UpdatedAt = time.Now()
	c.RecalculateTotal()
	return nil
}

func (c *ShoppingCart) RemoveItem(itemID id.UUID) error {
	updatedItems := make([]*ShoppingCartItem, 0, len(c.Items))
	removed := false
	for _, item := range c.Items {
		if item.ID == itemID {
			removed = true
			continue
		}
		updatedItems = append(updatedItems, item)
	}

	if !removed {
		return apperrors.NotFound(shoppingCartEntity+".remove_item: shopping cart item not found", nil)
	}

	c.Items = updatedItems
	c.UpdatedAt = time.Now()
	c.RecalculateTotal()
	return nil
}

func (c *ShoppingCart) FindItem(itemID id.UUID) *ShoppingCartItem {
	for _, item := range c.Items {
		if item.ID == itemID {
			return item
		}
	}
	return nil
}

func (c *ShoppingCart) RecalculateTotal() {
	var total float64
	for _, item := range c.Items {
		total += item.Subtotal()
	}
	c.TotalBTC = total
	c.UpdatedAt = time.Now()
}
