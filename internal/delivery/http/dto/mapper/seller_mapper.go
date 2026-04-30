package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dseller"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type SellerMapper struct {
}

func NewSellerMapper() *SellerMapper {
	return &SellerMapper{}
}

func (m *SellerMapper) SellerToResponse(seller *entity.Seller) *dseller.SellerResponse {
	return &dseller.SellerResponse{
		ID:          seller.ID.String(),
		Bio:         seller.Bio,
		DisplayName: seller.DisplayName,
	}
}

func (m *SellerMapper) CreateToSeller(dto *dseller.CreateSeller) *entity.Seller {
	return &entity.Seller{
		Bio:         dto.Bio,
		DisplayName: dto.DisplayName,
	}
}

func (m *SellerMapper) UpdateToSeller(dto *dseller.UpdateSeller) *entity.Seller {
	return &entity.Seller{
		Bio:         dto.Bio,
		DisplayName: dto.DisplayName,
	}
}
