package repository

import (
	"errors"
	"flat_bot/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id int64) (model.User, error)
	FindAll() ([]model.User, error)
	Create(user model.User) (model.User, error)
	ExistsByID(id int64) (bool, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepositoryImpl{db: db}
}

func (r UserRepositoryImpl) FindByID(id int64) (model.User, error) {
	var user model.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (r UserRepositoryImpl) FindAll() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r UserRepositoryImpl) Create(user model.User) (model.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r UserRepositoryImpl) ExistsByID(id int64) (bool, error) {
	var user model.User
	if err := r.db.Select("id").Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
