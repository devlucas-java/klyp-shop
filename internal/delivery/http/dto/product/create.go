package product

import "github.com/devlucas-java/klyp-shop/internal/domain/errors"

type CreateProduct struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PriceBTC    float64  `json:"price_btc"`
	Stock       int      `json:"stock"`
	Categories  []string `json:"categories"`
}

func (r *CreateProduct) Validate() error {
	if len(r.Name) < 2 || len(r.Name) > 200 {
		return errors.ErrBadRequest("name must be between 2 and 200 characters", nil)
	}
	if len(r.Description) > 5000 {
		return errors.ErrBadRequest("description must not exceed 5000 characters", nil)
	}
	if r.PriceBTC <= 0 {
		return errors.ErrBadRequest("price_btc must be greater than 0", nil)
	}
	if r.Stock < 0 {
		return errors.ErrBadRequest("stock must be greater than or equal to 0", nil)
	}
	return nil
}
