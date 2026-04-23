package daddress

type UpdateAddressRequest struct {
	Street   string
	City     string
	State    string
	Country  string
	PostCode string
	Number   int32
}
