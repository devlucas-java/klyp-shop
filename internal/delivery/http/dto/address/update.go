package daddress

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

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
		return apperrors.BadRequest("at least one field must be provided", nil)
	}
	if r.Street != "" && (len(r.Street) < 3 || len(r.Street) > 200) {
		return apperrors.BadRequest("street must be between 3 and 200 characters", nil)
	}
	if r.City != "" && len(r.City) > 100 {
		return apperrors.BadRequest("city must not exceed 100 characters", nil)
	}
	if r.State != "" && len(r.State) > 100 {
		return apperrors.BadRequest("state must not exceed 100 characters", nil)
	}
	if r.Country != "" && len(r.Country) > 100 {
		return apperrors.BadRequest("country must not exceed 100 characters", nil)
	}
	if r.PostCode != "" && len(r.PostCode) > 20 {
		return apperrors.BadRequest("postCode must not exceed 20 characters", nil)
	}
	if r.Number < 0 {
		return apperrors.BadRequest("number must be greater than or equal to 0", nil)
	}
	return nil
}
