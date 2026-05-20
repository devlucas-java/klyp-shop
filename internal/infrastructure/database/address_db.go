package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type AddressDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAddressDB(db *gorm.DB, log *logger.Logger) repository.AddressRepository {
	return &AddressDB{db: db, log: log}
}

func (a *AddressDB) Create(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Create(address).Error; err != nil {
		return nil, errors.HandlePgError(a.log, err, "failed to create address")
	}
	return address, nil
}

func (a *AddressDB) Save(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Where("id = ?", address.ID).Save(address).Error; err != nil {
		return nil, errors.HandlePgError(a.log, err, "failed to save address")
	}
	return address, nil
}

func (a *AddressDB) Updates(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Model(address).Where("id = ?", address.ID).Updates(address).Error; err != nil {
		return nil, errors.HandlePgError(a.log, err, "failed to update address")
	}
	return address, nil
}

func (a *AddressDB) Update(address *entity.Address) (*entity.Address, error) {
	return a.Save(address)
}

func (a *AddressDB) FindByID(addressID id.UUID) (*entity.Address, error) {
	var address entity.Address
	err := a.db.WithContext(context.Background()).First(&address, "id = ?", addressID).Error
	if err != nil {
		return nil, errors.HandlePgError(a.log, err, "failed to find address")
	}
	return &address, nil
}

func (a *AddressDB) FindByUser(userID id.UUID) ([]*entity.Address, error) {
	var addresses []*entity.Address
	if err := a.db.WithContext(context.Background()).Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, errors.HandlePgError(a.log, err, "failed to find addresses")
	}
	return addresses, nil
}

func (a *AddressDB) DeleteByID(addressID id.UUID) error {
	if err := a.db.WithContext(context.Background()).Delete(&entity.Address{}, "id = ?", addressID).Error; err != nil {
		return errors.HandlePgError(a.log, err, "failed to delete address")
	}
	return nil
}
