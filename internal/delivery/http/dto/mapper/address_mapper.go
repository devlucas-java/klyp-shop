package mapper

import (
	addressDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/address"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type AddressMapper struct{}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{}
}

func (m *AddressMapper) ToResponse(addr *entity.Address) *addressDTO.AddressResponse {
	return &addressDTO.AddressResponse{
		ID:       addr.ID.String(),
		Street:   addr.Street,
		City:     addr.City,
		State:    addr.State,
		Country:  addr.Country,
		PostCode: addr.Postcode,
		Number:   addr.Number,
	}
}
