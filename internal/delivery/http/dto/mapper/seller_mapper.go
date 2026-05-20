package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dseller"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type SellerMapper struct{}

func NewSellerMapper() *SellerMapper {
	return &SellerMapper{}
}

func (m *SellerMapper) ToResponse(seller *entity.Seller) *dseller.SellerResponse {
	return &dseller.SellerResponse{
		ID:          seller.ID.String(),
		Bio:         seller.Bio,
		DisplayName: seller.DisplayName,
	}
}
