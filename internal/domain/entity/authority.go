package entity

type Authority struct {
	BaseModel
	Name string
}

func NewAuthority(name string) *Authority {
	return &Authority{Name: name}
}
