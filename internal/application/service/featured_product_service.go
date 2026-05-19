package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dproduct"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type FeaturedProductService struct {
	log                *logger.Logger
	featuredRepository repository.FeaturedProductRepository
	productRepository  repository.ProductRepository
	userRepository     repository.UserRepository
}

func NewFeaturedProductService(
	log *logger.Logger,
	featuredRepository repository.FeaturedProductRepository,
	productRepository repository.ProductRepository,
	userRepository repository.UserRepository,
) *FeaturedProductService {
	return &FeaturedProductService{
		log:                log,
		featuredRepository: featuredRepository,
		productRepository:  productRepository,
		userRepository:     userRepository,
	}
}

func (s *FeaturedProductService) AddFeatured(auth *entity.User, req *dproduct.AddFeaturedRequest) (*dproduct.FeaturedProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	productID, err := id.Parse(req.ProductID)
	if err != nil {
		return nil, errors.ErrInvalidUUID(err)
	}

	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, errors.ErrNotFound("Product", err)
	}

	if !product.IsOwnedBy(user.Seller.ID) {
		return nil, errors.ErrForbidden(nil)
	}

	count, err := s.featuredRepository.CountBySellerID(user.Seller.ID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to count featured products", err)
	}
	if count >= entity.MaxFeaturedProducts {
		return nil, errors.ErrUnprocessable("maximum of 10 featured products reached", nil)
	}

	existing, _ := s.featuredRepository.FindBySellerIDAndProductID(user.Seller.ID, productID)
	if existing != nil {
		return nil, errors.ErrConflict("FeaturedProduct", nil)
	}

	featured, err := entity.NewFeaturedProduct(user.Seller.ID, productID, req.Position)
	if err != nil {
		return nil, err
	}

	saved, err := s.featuredRepository.Add(featured)
	if err != nil {
		return nil, errors.ErrDatabase("failed to add featured product", err)
	}

	return toFeaturedResponse(saved, product), nil
}

func (s *FeaturedProductService) RemoveFeatured(auth *entity.User, productID id.UUID) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return errors.ErrNotFound("User", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return err
	}

	return s.featuredRepository.Remove(user.Seller.ID, productID)
}

func (s *FeaturedProductService) UpdatePosition(auth *entity.User, productID id.UUID, req *dproduct.UpdateFeaturedPositionRequest) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return errors.ErrNotFound("User", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return err
	}

	_, err = s.featuredRepository.FindBySellerIDAndProductID(user.Seller.ID, productID)
	if err != nil {
		return errors.ErrNotFound("FeaturedProduct", err)
	}

	return s.featuredRepository.UpdatePosition(user.Seller.ID, productID, req.Position)
}

func (s *FeaturedProductService) GetFeaturedBySeller(sellerID id.UUID) ([]*dproduct.FeaturedProductResponse, error) {
	featured, err := s.featuredRepository.FindBySellerID(sellerID)
	if err != nil {
		return nil, errors.ErrDatabase("failed to fetch featured products", err)
	}

	result := make([]*dproduct.FeaturedProductResponse, 0, len(featured))
	for _, f := range featured {
		result = append(result, toFeaturedResponse(f, &f.Product))
	}
	return result, nil
}

func (s *FeaturedProductService) GetMyFeatured(auth *entity.User) ([]*dproduct.FeaturedProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	return s.GetFeaturedBySeller(user.Seller.ID)
}

func toFeaturedResponse(f *entity.FeaturedProduct, p *entity.Product) *dproduct.FeaturedProductResponse {
	return &dproduct.FeaturedProductResponse{
		ID:       f.ID.String(),
		Position: f.Position,
		Product: dproduct.ProductResponse{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			PriceBTC:    p.PriceBTC,
			Stock:       p.Stock,
			SellerID:    p.SellerID.String(),
			Categories:  p.Categories,
		},
	}
}
