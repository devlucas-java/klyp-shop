package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/product"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const featuredProductServiceTrace = "featured_product_service.FeaturedProductService"

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

func (s *FeaturedProductService) AddFeatured(auth *entity.User, req *product.AddFeaturedRequest) (*product.FeaturedProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(featuredProductServiceTrace+".add_featured: user not found", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	productID, err := id.Parse(req.ProductID)
	if err != nil {
		return nil, apperrors.InvalidUUID(featuredProductServiceTrace+".add_featured: invalid product id", err)
	}

	prod, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, apperrors.NotFound(featuredProductServiceTrace+".add_featured: product not found", err)
	}

	if !prod.IsOwnedBy(user.Seller.ID) {
		return nil, apperrors.Forbidden(featuredProductServiceTrace+".add_featured: product does not belong to seller", nil)
	}

	count, err := s.featuredRepository.CountBySellerID(user.Seller.ID)
	if err != nil {
		return nil, apperrors.Database(featuredProductServiceTrace+".add_featured: failed to count featured products", err)
	}
	if count >= entity.MaxFeaturedProducts {
		return nil, apperrors.Unprocessable(featuredProductServiceTrace+".add_featured: maximum of 10 featured products reached", nil)
	}

	existing, _ := s.featuredRepository.FindBySellerIDAndProductID(user.Seller.ID, productID)
	if existing != nil {
		return nil, apperrors.Conflict(featuredProductServiceTrace+".add_featured: product is already featured", nil)
	}

	featured, err := entity.NewFeaturedProduct(user.Seller.ID, productID, req.Position)
	if err != nil {
		return nil, err
	}

	saved, err := s.featuredRepository.Add(featured)
	if err != nil {
		return nil, apperrors.Database(featuredProductServiceTrace+".add_featured: failed to add featured product", err)
	}

	return toFeaturedResponse(saved, prod), nil
}

func (s *FeaturedProductService) RemoveFeatured(auth *entity.User, productID id.UUID) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return apperrors.NotFound(featuredProductServiceTrace+".remove_featured: user not found", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return err
	}

	return s.featuredRepository.Remove(user.Seller.ID, productID)
}

func (s *FeaturedProductService) UpdatePosition(auth *entity.User, productID id.UUID, req *product.UpdateFeaturedPositionRequest) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return apperrors.NotFound(featuredProductServiceTrace+".update_position: user not found", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return err
	}

	_, err = s.featuredRepository.FindBySellerIDAndProductID(user.Seller.ID, productID)
	if err != nil {
		return apperrors.NotFound(featuredProductServiceTrace+".update_position: featured product not found", err)
	}

	return s.featuredRepository.UpdatePosition(user.Seller.ID, productID, req.Position)
}

func (s *FeaturedProductService) GetAllFeatured() ([]*product.FeaturedProductResponse, error) {
	featured, err := s.featuredRepository.FindAll()
	if err != nil {
		return nil, apperrors.Database(featuredProductServiceTrace+".get_all_featured: failed to fetch featured products", err)
	}

	result := make([]*product.FeaturedProductResponse, 0, len(featured))
	for _, f := range featured {
		result = append(result, toFeaturedResponse(f, &f.Product))
	}
	return result, nil
}

func (s *FeaturedProductService) GetFeaturedBySeller(sellerID id.UUID) ([]*product.FeaturedProductResponse, error) {
	featured, err := s.featuredRepository.FindBySellerID(sellerID)
	if err != nil {
		return nil, apperrors.Database(featuredProductServiceTrace+".get_featured_by_seller: failed to fetch featured products", err)
	}

	result := make([]*product.FeaturedProductResponse, 0, len(featured))
	for _, f := range featured {
		result = append(result, toFeaturedResponse(f, &f.Product))
	}
	return result, nil
}

func (s *FeaturedProductService) GetMyFeatured(auth *entity.User) ([]*product.FeaturedProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(featuredProductServiceTrace+".get_my_featured: user not found", err)
	}
	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	return s.GetFeaturedBySeller(user.Seller.ID)
}

func toFeaturedResponse(f *entity.FeaturedProduct, p *entity.Product) *product.FeaturedProductResponse {
	return &product.FeaturedProductResponse{
		ID:       f.ID.String(),
		Position: f.Position,
		Product: product.ProductResponse{
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
