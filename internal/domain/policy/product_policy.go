package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ProductPolicy struct{}

func NewProductPolicy() *ProductPolicy {
	return &ProductPolicy{}
}

func (p *ProductPolicy) CanManage(product *entity.Product, sellerID id.UUID) error {
	if !product.IsOwnedBy(sellerID) {
		return errors.ErrForbidden(nil)
	}
	return nil
}
