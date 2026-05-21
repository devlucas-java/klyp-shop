package service

import (
	addressDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/address"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const addressService = "address_service"

type AddressService struct {
	addressRepository repository.AddressRepository
	userRepository    repository.UserRepository
	log               *logger.Logger
	mapper            *mapper.AddressMapper
	addressPolicy     *policy.AddressPolicy
}

func NewAddressService(
	addressRepository repository.AddressRepository,
	userRepository repository.UserRepository,
	log *logger.Logger,
	mapper *mapper.AddressMapper,
	ap *policy.AddressPolicy,
) *AddressService {
	return &AddressService{
		addressRepository: addressRepository,
		userRepository:    userRepository,
		log:               log,
		mapper:            mapper,
		addressPolicy:     ap,
	}
}

func (s *AddressService) CreateAddress(auth *entity.User, req *addressDTO.CreateAddressRequest) (*addressDTO.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(addressService+".create_address: user not found", err)
	}

	existing, err := s.addressRepository.FindByUser(user.ID)
	if err != nil {
		return nil, apperrors.Database(addressService+".create_address: failed to fetch addresses", err)
	}

	if err := s.addressPolicy.CanCreate(existing); err != nil {
		return nil, err
	}

	address := entity.NewAddress(user.ID, req.Street, req.City, req.State, req.Country, req.PostCode, req.Number)

	saved, err := s.addressRepository.Create(address)
	if err != nil {
		return nil, apperrors.Database(addressService+".create_address: failed to create address", err)
	}

	return s.mapper.ToResponse(saved), nil
}

func (s *AddressService) GetAddresses(auth *entity.User) ([]*addressDTO.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(addressService+".get_addresses: user not found", err)
	}

	addrs, err := s.addressRepository.FindByUser(user.ID)
	if err != nil {
		return nil, apperrors.Database(addressService+".get_addresses: failed to fetch addresses", err)
	}

	responses := make([]*addressDTO.AddressResponse, len(addrs))
	for i, addr := range addrs {
		responses[i] = s.mapper.ToResponse(addr)
	}

	return responses, nil
}

func (s *AddressService) UpdateAddress(auth *entity.User, req *addressDTO.UpdateAddressRequest, addrID id.UUID) (*addressDTO.AddressResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(addressService+".update_address: user not found", err)
	}

	addr, err := s.addressRepository.FindByID(addrID)
	if err != nil {
		return nil, apperrors.NotFound(addressService+".update_address: address not found", err)
	}

	if err := s.addressPolicy.CanModify(addr, user.ID); err != nil {
		return nil, err
	}

	if req.Street != "" {
		addr.Street = req.Street
	}
	if req.City != "" {
		addr.City = req.City
	}
	if req.State != "" {
		addr.State = req.State
	}
	if req.Country != "" {
		addr.Country = req.Country
	}
	if req.PostCode != "" {
		addr.Postcode = req.PostCode
	}
	if req.Number > 0 {
		addr.Number = req.Number
	}

	saved, err := s.addressRepository.Update(addr)
	if err != nil {
		return nil, apperrors.Database(addressService+".update_address: failed to update address", err)
	}

	return s.mapper.ToResponse(saved), nil
}

func (s *AddressService) DeleteAddress(auth *entity.User, addrID id.UUID) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return apperrors.NotFound(addressService+".delete_address: user not found", err)
	}

	addr, err := s.addressRepository.FindByID(addrID)
	if err != nil {
		return apperrors.NotFound(addressService+".delete_address: address not found", err)
	}

	if err := s.addressPolicy.CanModify(addr, user.ID); err != nil {
		return err
	}

	if err := s.addressRepository.DeleteByID(addr.ID); err != nil {
		return apperrors.Database(addressService+".delete_address: failed to delete address", err)
	}

	return nil
}
