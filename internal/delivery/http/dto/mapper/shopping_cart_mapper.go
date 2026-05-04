package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dcart"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type ShoppingCartMapper struct{}

func NewShoppingCartMapper() *ShoppingCartMapper {
	return &ShoppingCartMapper{}
}

func (m *ShoppingCartMapper) ShoppingCartToResponse(cart *entity.ShoppingCart) *dcart.ShoppingCartResponse {
	if cart == nil {
		return nil
	}

	items := make([]dcart.ShoppingCartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = dcart.ShoppingCartItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			PriceBTC:  item.PriceBTC,
			Subtotal:  item.Subtotal(),
		}
	}

	return &dcart.ShoppingCartResponse{
		ID:        cart.ID.String(),
		UserID:    cart.UserID.String(),
		TotalBTC:  cart.TotalBTC,
		Items:     items,
		CreatedAt: cart.CreatedAt.String(),
		UpdatedAt: cart.UpdatedAt.String(),
	}
}
