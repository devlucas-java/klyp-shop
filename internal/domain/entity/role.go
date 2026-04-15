package entity

type Role struct {
	BaseModel

	Name string

	Authorities []Authority `gorm:"many2many:role_authorities"`
}

func NewRole(name string) *Role {
	return &Role{
		Name: name,
	}
}
