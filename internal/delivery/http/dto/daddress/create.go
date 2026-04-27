package daddress

type CreateAddressRequest struct {
	Street   string `json:"street" validate:"required,min=3,max=100"`
	City     string `json:"city" validate:"required"`
	State    string `json:"state" validate:"required"`
	Country  string `json:"country" validate:"required,iso3166_1_alpha2"`
	PostCode string `json:"postCode" validate:"required"`
	Number   int32  `json:"number" validate:"required,gt=0"`
}
