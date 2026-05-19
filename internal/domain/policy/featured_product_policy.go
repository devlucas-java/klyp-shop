package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

// FeaturedProductPolicy contém as regras de negócio para produtos em destaque.
type FeaturedProductPolicy struct{}

func NewFeaturedProductPolicy() *FeaturedProductPolicy {
	return &FeaturedProductPolicy{}
}

// CanAdd verifica se o seller pode adicionar mais um produto em destaque.
func (p *FeaturedProductPolicy) CanAdd(currentCount int64) error {
	if currentCount >= entity.MaxFeaturedProducts {
		return errors.ErrUnprocessable("maximum of 10 featured products reached", nil)
	}
	return nil
}

// CanManage verifica se o produto pertence ao seller que quer gerenciá-lo.
func (p *FeaturedProductPolicy) CanManage(product *entity.Product, sellerID id.UUID) error {
	if !product.IsOwnedBy(sellerID) {
		return errors.ErrForbidden(nil)
	}
	return nil
}
