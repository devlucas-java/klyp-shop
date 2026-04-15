package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	b.ID = id.NewUUID()
	return nil
}
