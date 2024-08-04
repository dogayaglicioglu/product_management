package repository

import (
	"fmt"
	"product_management/models"

	"gorm.io/gorm"
)

type DbRepository interface {
	GetUserById(id string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}

type dbRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) DbRepository {
	return &dbRepo{db: db} //dbRepo should implement DbRepository
}

func (d *dbRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if result := d.db.Find(&users); result.Error != nil {
		return nil, result.Error
		fmt.Errorf("Error in fetching users from db %v", result.Error)

	}
	return users, nil
}

func (d *dbRepo) GetUserById(id string) (*models.User, error) {
	var user models.User
	if result := d.db.First(&user, "username = ?", id); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
