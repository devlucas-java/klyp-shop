package database

import (
	"errors"
	"fmt"

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
	if err := a.db.Create(address).Error; err != nil {
		a.log.Errorf("AddressDB.Create: %v", err)
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	return address, nil
}

func (a *AddressDB) Save(address *entity.Address) (*entity.Address, error) {
	if err := a.db.Where("id = ?", address.ID).Save(address).Error; err != nil {
		a.log.Errorf("AddressDB.Save %s: %v", address.ID, err)
		return nil, fmt.Errorf("failed to save address: %w", err)
	}
	return address, nil
}

func (a *AddressDB) Updates(address *entity.Address) (*entity.Address, error) {
	if err := a.db.Model(address).Where("id = ?", address.ID).Updates(address).Error; err != nil {
		a.log.Errorf("AddressDB.Updates %s: %v", address.ID, err)
		return nil, fmt.Errorf("failed to update address: %w", err)
	}
	return address, nil
}

func (a *AddressDB) Update(address *entity.Address) (*entity.Address, error) {
	return a.Save(address)
}

func (a *AddressDB) FindByID(addressID id.UUID) (*entity.Address, error) {
	var address entity.Address
	err := a.db.First(&address, "id = ?", addressID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Address", err)
		}
		a.log.Errorf("AddressDB.FindByID %s: %v", addressID, err)
		return nil, fmt.Errorf("failed to find address: %w", err)
	}
	return &address, nil
}

func (a *AddressDB) FindByUser(userID id.UUID) ([]*entity.Address, error) {
	var addresses []*entity.Address
	if err := a.db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		a.log.Errorf("AddressDB.FindByUser %s: %v", userID, err)
		return nil, fmt.Errorf("failed to find addresses: %w", err)
	}
	return addresses, nil
}

func (a *AddressDB) DeleteByID(addressID id.UUID) error {
	if err := a.db.Delete(&entity.Address{}, "id = ?", addressID).Error; err != nil {
		a.log.Errorf("AddressDB.DeleteByID %s: %v", addressID, err)
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}
