package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type AddressMapper struct {
}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{}
}

func (m *AddressMapper) AddressToAddressResponse(addr *entity.Address) *daddress.AddressResponse {

	return &daddress.AddressResponse{
		Street:   addr.Street,
		City:     addr.City,
		State:    addr.State,
		Country:  addr.Country,
		PostCode: addr.Postcode,
		Number:   addr.Number,
	}
}

func (m *AddressMapper) CreateAddressRequestToAddress(dto *daddress.CreateAddressRequest) *entity.Address {

	return &entity.Address{
		Street:   dto.Street,
		City:     dto.City,
		State:    dto.State,
		Country:  dto.Country,
		Postcode: dto.PostCode,
		Number:   dto.Number,
	}
}

func (m *AddressMapper) AddressDTOToAddress(dto *daddress.AddressResponse) *entity.Address {

	return &entity.Address{
		Street:   dto.Street,
		City:     dto.City,
		State:    dto.State,
		Country:  dto.Country,
		Postcode: dto.PostCode,
		Number:   dto.Number,
	}
}
