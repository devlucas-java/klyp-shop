package daddress

import "github.com/devlucas-java/klyp-shop/internal/domain/errors"

type CreateAddressRequest struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	PostCode string `json:"postCode"`
	Number   int32  `json:"number"`
}

func (r *CreateAddressRequest) Validate() error {
	if len(r.Street) < 3 || len(r.Street) > 200 {
		return errors.ErrBadRequest("street must be between 3 and 200 characters", nil)
	}
	if r.City == "" || len(r.City) > 100 {
		return errors.ErrBadRequest("city is required and must not exceed 100 characters", nil)
	}
	if r.State == "" || len(r.State) > 100 {
		return errors.ErrBadRequest("state is required and must not exceed 100 characters", nil)
	}
	if r.Country == "" || len(r.Country) > 100 {
		return errors.ErrBadRequest("country is required and must not exceed 100 characters", nil)
	}
	if r.PostCode == "" || len(r.PostCode) > 20 {
		return errors.ErrBadRequest("postCode is required and must not exceed 20 characters", nil)
	}
	if r.Number <= 0 {
		return errors.ErrBadRequest("number must be greater than 0", nil)
	}
	return nil
}
