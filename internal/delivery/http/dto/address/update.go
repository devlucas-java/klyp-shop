package daddress

import "github.com/devlucas-java/klyp-shop/internal/domain/errors"

type UpdateAddressRequest struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	PostCode string `json:"postCode"`
	Number   int32  `json:"number"`
}

func (r *UpdateAddressRequest) Validate() error {
	if r.Street == "" && r.City == "" && r.State == "" && r.Country == "" && r.PostCode == "" && r.Number == 0 {
		return errors.ErrBadRequest("at least one field must be provided", nil)
	}
	if r.Street != "" && (len(r.Street) < 3 || len(r.Street) > 200) {
		return errors.ErrBadRequest("street must be between 3 and 200 characters", nil)
	}
	if r.City != "" && len(r.City) > 100 {
		return errors.ErrBadRequest("city must not exceed 100 characters", nil)
	}
	if r.State != "" && len(r.State) > 100 {
		return errors.ErrBadRequest("state must not exceed 100 characters", nil)
	}
	if r.Country != "" && len(r.Country) > 100 {
		return errors.ErrBadRequest("country must not exceed 100 characters", nil)
	}
	if r.PostCode != "" && len(r.PostCode) > 20 {
		return errors.ErrBadRequest("postCode must not exceed 20 characters", nil)
	}
	if r.Number < 0 {
		return errors.ErrBadRequest("number must be greater than or equal to 0", nil)
	}
	return nil
}
