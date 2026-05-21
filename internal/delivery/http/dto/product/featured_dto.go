package product

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type AddFeaturedRequest struct {
	ProductID string `json:"product_id"`
	Position  int    `json:"position"`
}

func (r *AddFeaturedRequest) Validate() error {
	if r.ProductID == "" {
		return apperrors.BadRequest("product_id is required", nil)
	}
	if r.Position < 1 || r.Position > 10 {
		return apperrors.BadRequest("position must be between 1 and 10", nil)
	}
	return nil
}

type UpdateFeaturedPositionRequest struct {
	Position int `json:"position"`
}

func (r *UpdateFeaturedPositionRequest) Validate() error {
	if r.Position < 1 || r.Position > 10 {
		return apperrors.BadRequest("position must be between 1 and 10", nil)
	}
	return nil
}

type FeaturedProductResponse struct {
	ID       string          `json:"id"`
	Position int             `json:"position"`
	Product  ProductResponse `json:"product"`
}
