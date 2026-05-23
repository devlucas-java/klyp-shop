package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ProductPolicy struct{}

func NewProductPolicy() *ProductPolicy {
	return &ProductPolicy{}
}

func (p *ProductPolicy) CanManage(product *entity.Product, sellerID id.UUID) error {
	if !product.IsOwnedBy(sellerID) {
		return apperrors.Forbidden(nil)
	}
	return nil
}
