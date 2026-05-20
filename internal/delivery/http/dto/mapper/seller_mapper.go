package mapper

import (
	sellerDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/seller"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type SellerMapper struct{}

func NewSellerMapper() *SellerMapper {
	return &SellerMapper{}
}

func (m *SellerMapper) ToResponse(seller *entity.Seller) *sellerDTO.SellerResponse {
	return &sellerDTO.SellerResponse{
		ID:          seller.ID.String(),
		Bio:         seller.Bio,
		DisplayName: seller.DisplayName,
	}
}
