package service

import (
	"github.com/kecci/goscription/internal/repository/postgres"
	"github.com/kecci/goscription/models"
)

type (
	AddressService interface {
		Insert(address models.Address) error
		GetAddressAll() ([]models.Address, error)
	}

	AddressServiceImpl struct {
		addressRepo postgres.AddressRepository
	}
)

func NewAddressService(addressRepo postgres.AddressRepository) AddressService {
	return &AddressServiceImpl{
		addressRepo: addressRepo,
	}
}

func (a *AddressServiceImpl) Insert(address models.Address) error {
	return a.addressRepo.Insert(address)
}

func (a *AddressServiceImpl) GetAddressAll() ([]models.Address, error) {
	return a.addressRepo.GetAddressAll()
}
