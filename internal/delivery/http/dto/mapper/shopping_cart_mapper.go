package mapper

import (
	cartDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/cart"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type ShoppingCartMapper struct{}

func NewShoppingCartMapper() *ShoppingCartMapper {
	return &ShoppingCartMapper{}
}

func (m *ShoppingCartMapper) ShoppingCartToResponse(cart *entity.ShoppingCart) *cartDTO.ShoppingCartResponse {
	if cart == nil {
		return nil
	}

	items := make([]cartDTO.ShoppingCartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = cartDTO.ShoppingCartItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			PriceBTC:  item.PriceBTC,
			Subtotal:  item.Subtotal(),
		}
	}

	return &cartDTO.ShoppingCartResponse{
		ID:        cart.ID.String(),
		UserID:    cart.UserID.String(),
		TotalBTC:  cart.TotalBTC,
		Items:     items,
		CreatedAt: cart.CreatedAt.String(),
		UpdatedAt: cart.UpdatedAt.String(),
	}
}
