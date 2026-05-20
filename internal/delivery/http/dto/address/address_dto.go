package daddress

type AddressResponse struct {
	ID       string `json:"id"`
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	PostCode string `json:"postCode"`
	Number   int32  `json:"number"`
}
