package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type AddressService struct {
	addressRepository repository.AddressRepository
	log               *logger.Logger
	mapper            *mapper.AddressMapper
	userRepository    repository.UserRepository
}

func NewAddressService(addressRepository repository.AddressRepository, log *logger.Logger, mapper *mapper.AddressMapper, userRepository repository.UserRepository) *AddressService {
	return &AddressService{addressRepository: addressRepository, log: log, mapper: mapper, userRepository: userRepository}
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

	if len(addrs) >= 3 {
		return nil, errors.ErrUnprocessable("maximum number of address (3) reached", err)
	}
	_, err = s.addressRepository.Create(address)
	if err != nil {
		return nil, err
	}

	return s.mapper.AddressToAddressResponse(address), nil
}

func (s *AddressService) GetAddresses(auth *entity.User) ([]*daddress.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, err
	}

	addrs, err := s.addressRepository.FindByUser(user.ID)
	if err != nil {
		return nil, err
	}

	var response []*daddress.AddressResponse
	for _, addr := range addrs {
		response = append(response, s.mapper.AddressToAddressResponse(addr))
	}

	return response, nil
}

func (s *AddressService) UpdateAddress(auth *entity.User, req *daddress.UpdateAddressRequest, id id.UUID) (*daddress.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, err
	}

	addr, err := s.addressRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if addr.UserID != user.ID {
		return nil, errors.ErrUnprocessable("you can not update this address", nil)
	}

	addr = s.mapper.UpdateAddressRequestToAddress(req)
	addr.UserID = user.ID
	addr.ID = id

	addr, err = s.addressRepository.Update(addr)
	if err != nil {
		return nil, err
	}

	return s.mapper.AddressToAddressResponse(addr), nil
}

func (s *AddressService) DeleteAddress(auth *entity.User, req id.UUID) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return err
	}

	addr, err := s.addressRepository.FindByID(req)
	if err != nil {
		return err
	}

	if addr.UserID != user.ID {
		return errors.ErrUnprocessable("you can not delete this address", nil)
	}

	err = s.addressRepository.DeleteByID(addr.ID)
	if err != nil {
		return err
	}

	return nil
}
