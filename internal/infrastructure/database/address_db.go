package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

const addressDB = "address_db"

type AddressDB struct {
	db *gorm.DB
}

func NewAddressDB(db *gorm.DB) repository.AddressRepository {
	return &AddressDB{db: db}
}

func (a *AddressDB) Create(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Create(address).Error; err != nil {
		return nil, apperrors.HandlePgError(addressDB+".create", err)
	}
	return address, nil
}

func (a *AddressDB) Save(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Where("id = ?", address.ID).Save(address).Error; err != nil {
		return nil, apperrors.HandlePgError(addressDB+".save", err)
	}
	return address, nil
}

func (a *AddressDB) Updates(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Model(address).Where("id = ?", address.ID).Updates(address).Error; err != nil {
		return nil, apperrors.HandlePgError(addressDB+".updates", err)
	}
	return address, nil
}

func (a *AddressDB) Update(address *entity.Address) (*entity.Address, error) {
	saved, err := a.Save(address)
	if err != nil {
		return nil, apperrors.HandlePgError(addressDB+".update", err)
	}
	return saved, nil
}

func (a *AddressDB) FindByID(addressID id.UUID) (*entity.Address, error) {
	var address entity.Address
	if err := a.db.WithContext(context.Background()).First(&address, "id = ?", addressID).Error; err != nil {
		return nil, apperrors.HandlePgError(addressDB+".find_by_id", err)
	}
	return &address, nil
}

func (a *AddressDB) FindByUser(userID id.UUID) ([]*entity.Address, error) {
	var addresses []*entity.Address
	if err := a.db.WithContext(context.Background()).Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, apperrors.HandlePgError(addressDB+".find_by_user", err)
	}
	return addresses, nil
}

func (a *AddressDB) DeleteByID(addressID id.UUID) error {
	if err := a.db.WithContext(context.Background()).Delete(&entity.Address{}, "id = ?", addressID).Error; err != nil {
		return apperrors.HandlePgError(addressDB+".delete_by_id", err)
	}
	return nil
}
