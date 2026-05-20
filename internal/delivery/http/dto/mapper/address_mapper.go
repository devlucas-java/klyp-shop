package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type AddressMapper struct{}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{}
}

func (m *AddressMapper) ToResponse(addr *entity.Address) *daddress.AddressResponse {
	return &daddress.AddressResponse{
		ID:       addr.ID.String(),
		Street:   addr.Street,
		City:     addr.City,
		State:    addr.State,
		Country:  addr.Country,
		PostCode: addr.Postcode,
		Number:   addr.Number,
	}
}
