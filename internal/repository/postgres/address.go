package postgres

import (
	"fmt"

	"github.com/kecci/goscription/internal/library/db"
	"github.com/kecci/goscription/models"
	"gorm.io/gorm"
)

type (
	AddressRepository interface {
		Insert(address models.Address) (err error)
		GetAddressAll() (addresses []models.Address, err error)
	}

	AddressRepositoryImpl struct {
		DB *gorm.DB
	}
)

func NewAddressRepository(db db.Database) AddressRepository {
	if db.Postgres == nil {
		fmt.Println("db postgress nil")
	}
	return &AddressRepositoryImpl{DB: db.Postgres}
}

func (r *AddressRepositoryImpl) Insert(address models.Address) (err error) {
	tx := r.DB.Table("address").Create(&address).Scan(&address)
	if err = tx.Error; err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (r *AddressRepositoryImpl) GetAddressAll() (addresses []models.Address, err error) {
	tx := r.DB.Raw("SELECT * FROM address").Scan(&addresses)
	if tx.Error != nil {
		err = tx.Error
	}
	return
}
