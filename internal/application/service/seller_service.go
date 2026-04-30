package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dseller"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type SellerService struct {
	log              *logger.Logger
	userRepository   repository.UserRepository
	sellerRepository repository.SellerRepository
	sellerMapper     *mapper.SellerMapper
}

func NewSellerService(
	log *logger.Logger,
	userRepository repository.UserRepository,
	sellerRepository repository.SellerRepository,
	sellerMapper *mapper.SellerMapper,
) *SellerService {
	return &SellerService{
		log:              log,
		userRepository:   userRepository,
		sellerRepository: sellerRepository,
		sellerMapper:     sellerMapper,
	}
}

func (s *SellerService) CreateSeller(auth *entity.User, req *dseller.CreateSeller) (*dseller.SellerResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	if user.IsSeller {
		return nil, errors.ErrConflict("Seller", nil)
	}

	seller := s.sellerMapper.CreateToSeller(req)
	seller.UserID = user.ID

	saved, err := s.sellerRepository.Create(seller)
	if err != nil {
		return nil, errors.ErrDatabase("failed to create seller", err)
	}

	user.IsSeller = true
	if _, err = s.userRepository.Update(user); err != nil {
		return nil, errors.ErrDatabase("failed to update user", err)
	}

	return s.sellerMapper.SellerToResponse(saved), nil
}

func (s *SellerService) GetSellerByID(uuid id.UUID) (*dseller.SellerResponse, error) {
	seller, err := s.sellerRepository.FindByID(uuid)
	if err != nil {
		return nil, err
	}
	return s.sellerMapper.SellerToResponse(seller), nil
}

func (s *SellerService) UpdateSeller(auth *entity.User, req *dseller.UpdateSeller) (*dseller.SellerResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	if !user.IsSeller || user.Seller == nil {
		return nil, errors.ErrNotFound("Seller", nil)
	}

	patch := s.sellerMapper.UpdateToSeller(req)
	patch.ID = user.Seller.ID

	updated, err := s.sellerRepository.Updates(patch)
	if err != nil {
		return nil, errors.ErrDatabase("failed to update seller", err)
	}

	return s.sellerMapper.SellerToResponse(updated), nil
}

func (s *SellerService) DeleteSeller(auth *entity.User) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return errors.ErrNotFound("User", err)
	}

	if !user.IsSeller || user.Seller == nil {
		return errors.ErrNotFound("Seller", nil)
	}

	if err = s.sellerRepository.DeleteByID(user.Seller.ID); err != nil {
		return err
	}

	user.IsSeller = false
	user.Seller = nil
	if _, err = s.userRepository.Update(user); err != nil {
		return errors.ErrDatabase("failed to update user", err)
	}

	return nil
}
