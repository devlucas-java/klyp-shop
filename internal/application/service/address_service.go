package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type AddressService struct {
	addressRepository repository.AddressRepository
	log               *logger.Logger
	mapper            *mapper.AddressMapper
	userRepository    repository.UserRepository
}

func NewAddressService(addressRepository repository.AddressRepository, log *logger.Logger, mapper *mapper.AddressMapper) AddressService {
	return AddressService{addressRepository: addressRepository, log: log, mapper: mapper}
}

func (s *AddressService) CreateAddress(auth *entity.User, req *daddress.CreateAddressRequest) (*daddress.AddressResponse, error) {

	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, err
	}

	address := s.mapper.CreateAddressRequestToAddress(req)
	address.UserID = user.ID

	addrs, err := s.addressRepository.FindByUser(user.ID)
	if err != nil {
		return nil, err
	}

	if len(addrs) > 3 {
		return nil, errors.ErrUnprocessable("maximum number of address (3) reached", err)
	}
	_, err = s.addressRepository.Create(address)
	if err != nil {
		return nil, err
	}

	return s.mapper.AddressToAddressResponse(address), nil
}
