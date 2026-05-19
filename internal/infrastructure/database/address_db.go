package database

import (
	"context"
	"errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
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
		a.log.Errorf("AddressDB.Create: %v", err)
		return nil, domainErr.ErrDatabase("failed to create address", err)
	}
	return address, nil
}

func (a *AddressDB) Save(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Where("id = ?", address.ID).Save(address).Error; err != nil {
		a.log.Errorf("AddressDB.Save %s: %v", address.ID, err)
		return nil, domainErr.ErrDatabase("failed to save address", err)
	}
	return address, nil
}

func (a *AddressDB) Updates(address *entity.Address) (*entity.Address, error) {
	if err := a.db.WithContext(context.Background()).Model(address).Where("id = ?", address.ID).Updates(address).Error; err != nil {
		a.log.Errorf("AddressDB.Updates %s: %v", address.ID, err)
		return nil, domainErr.ErrDatabase("failed to update address", err)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Address", err)
		}
		a.log.Errorf("AddressDB.FindByID %s: %v", addressID, err)
		return nil, domainErr.ErrDatabase("failed to find address", err)
	}
	return &address, nil
}

func (a *AddressDB) FindByUser(userID id.UUID) ([]*entity.Address, error) {
	var addresses []*entity.Address
	if err := a.db.WithContext(context.Background()).Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		a.log.Errorf("AddressDB.FindByUser %s: %v", userID, err)
		return nil, domainErr.ErrDatabase("failed to find addresses", err)
	}
	return addresses, nil
}

func (a *AddressDB) DeleteByID(addressID id.UUID) error {
	if err := a.db.WithContext(context.Background()).Delete(&entity.Address{}, "id = ?", addressID).Error; err != nil {
		a.log.Errorf("AddressDB.DeleteByID %s: %v", addressID, err)
		return domainErr.ErrDatabase("failed to delete address", err)
	}
	return nil
}
