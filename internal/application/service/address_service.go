package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type AddressService struct {
	addressRepository repository.AddressRepository
	userRepository    repository.UserRepository
	log               *logger.Logger
	mapper            *mapper.AddressMapper
	addressPolicy     *policy.AddressPolicy
}

func NewAddressService(
	addressRepository repository.AddressRepository,
	log *logger.Logger,
	mapper *mapper.AddressMapper,
	userRepository repository.UserRepository,
) *AddressService {
	return &AddressService{
		addressRepository: addressRepository,
		userRepository:    userRepository,
		log:               log,
		mapper:            mapper,
		addressPolicy:     policy.NewAddressPolicy(),
	}
}

func (s *AddressService) CreateAddress(auth *entity.User, req *daddress.CreateAddressRequest) (*daddress.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	existing, err := s.addressRepository.FindByUser(user.ID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to fetch addresses", err)
	}

	if err := s.addressPolicy.CanCreate(existing); err != nil {
		return nil, err
	}

	address := s.mapper.CreateAddressRequestToAddress(req)
	address.UserID = user.ID

	saved, err := s.addressRepository.Create(address)
	if err != nil {
		return nil, errors.ErrDatabase("failed to create address", err)
	}

	return s.mapper.AddressToAddressResponse(saved), nil
}

func (s *AddressService) GetAddresses(auth *entity.User) ([]*daddress.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	addrs, err := s.addressRepository.FindByUser(user.ID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to fetch addresses", err)
	}

	responses := make([]*daddress.AddressResponse, len(addrs))
	for i, addr := range addrs {
		responses[i] = s.mapper.AddressToAddressResponse(addr)
	}

	return responses, nil
}

func (s *AddressService) UpdateAddress(auth *entity.User, req *daddress.UpdateAddressRequest, addrID id.UUID) (*daddress.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	addr, err := s.addressRepository.FindByID(addrID)
	if err != nil {
		return nil, errors.ErrNotFound("Address", err)
	}

	if err := s.addressPolicy.CanModify(addr, user.ID); err != nil {
		return nil, err
	}

	updated := s.mapper.UpdateAddressRequestToAddress(req)
	updated.ID = addrID
	updated.UserID = user.ID

	saved, err := s.addressRepository.Update(updated)
	if err != nil {
		return nil, errors.ErrDatabase("failed to update address", err)
	}

	return s.mapper.AddressToAddressResponse(saved), nil
}

func (s *AddressService) DeleteAddress(auth *entity.User, addrID id.UUID) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return errors.ErrNotFound("User", err)
	}

	addr, err := s.addressRepository.FindByID(addrID)
	if err != nil {
		return errors.ErrNotFound("Address", err)
	}

	if err := s.addressPolicy.CanModify(addr, user.ID); err != nil {
		return err
	}

	if err := s.addressRepository.DeleteByID(addr.ID); err != nil {
		return errors.ErrDatabase("failed to delete address", err)
	}

	return nil
}
