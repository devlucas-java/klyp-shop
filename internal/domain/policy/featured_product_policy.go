package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const featuredProductPolicy = "featured_product_policy.FeaturedProductPolicy"

type FeaturedProductPolicy struct{}

func NewFeaturedProductPolicy() *FeaturedProductPolicy {
	return &FeaturedProductPolicy{}
}

func (p *FeaturedProductPolicy) CanAdd(currentCount int64) error {
	if currentCount >= entity.MaxFeaturedProducts {
		return apperrors.Unprocessable(featuredProductPolicy+".can_add: maximum of 10 featured products reached", nil)
	}
	return nil
}

func (p *FeaturedProductPolicy) CanManage(product *entity.Product, sellerID id.UUID) error {
	if !product.IsOwnedBy(sellerID) {
		return apperrors.Forbidden(featuredProductPolicy+".can_manage: product does not belong to seller", nil)
	}
	return nil
}
