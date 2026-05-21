package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/seller"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const sellerServiceTrace = "seller_service.SellerService"

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

func (s *SellerService) CreateSeller(auth *entity.User, req *seller.CreateSeller) (*seller.SellerResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(sellerServiceTrace+".create_seller: user not found", err)
	}

	if err := user.MarkAsSeller(); err != nil {
		return nil, err
	}

	newSeller := entity.NewSeller(user.ID, req.DisplayName, req.Bio)

	saved, err := s.sellerRepository.Create(newSeller)
	if err != nil {
		return nil, apperrors.Database(sellerServiceTrace+".create_seller: failed to create seller", err)
	}

	user.ChangerToSeller()
	if _, err = s.userRepository.Update(user); err != nil {
		return nil, apperrors.Database(sellerServiceTrace+".create_seller: failed to update user", err)
	}

	return s.sellerMapper.ToResponse(saved), nil
}

func (s *SellerService) GetSellerByID(uuid id.UUID) (*seller.SellerResponse, error) {
	found, err := s.sellerRepository.FindByID(uuid)
	if err != nil {
		return nil, apperrors.NotFound(sellerServiceTrace+".get_seller_by_id: seller not found", err)
	}
	return s.sellerMapper.ToResponse(found), nil
}

func (s *SellerService) UpdateSeller(auth *entity.User, req *seller.UpdateSeller) (*seller.SellerResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(sellerServiceTrace+".update_seller: user not found", err)
	}

	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	user.Seller.UpdateInfo(req.DisplayName, req.Bio)

	updated, err := s.sellerRepository.Updates(user.Seller)
	if err != nil {
		return nil, apperrors.Database(sellerServiceTrace+".update_seller: failed to update seller", err)
	}

	return s.sellerMapper.ToResponse(updated), nil
}

func (s *SellerService) DeleteSeller(auth *entity.User) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return apperrors.NotFound(sellerServiceTrace+".delete_seller: user not found", err)
	}

	if err := user.EnsureSeller(); err != nil {
		return err
	}

	if err = s.sellerRepository.DeleteByID(user.Seller.ID); err != nil {
		return apperrors.Database(sellerServiceTrace+".delete_seller: failed to delete seller", err)
	}

	user.Seller = nil
	if err := user.UnmarkAsSeller(); err != nil {
		return err
	}
	user.ChangerToUser()

	if _, err = s.userRepository.Update(user); err != nil {
		return apperrors.Database(sellerServiceTrace+".delete_seller: failed to update user", err)
	}

	return nil
}
