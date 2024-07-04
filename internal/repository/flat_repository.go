package repository

import (
	"flat_bot/internal/model"
	"gorm.io/gorm"
)

type FlatRepository interface {
	FindByID(id string) (model.Flat, error)
	FindAll() ([]model.Flat, error)
	Create(flat model.Flat) (model.Flat, error)
	ExistsByID(id string) (bool, error)
}

type FlatRepositoryImpl struct {
	db *gorm.DB
}

func NewFlatRepository(db *gorm.DB) FlatRepository {
	return FlatRepositoryImpl{db: db}
}

func (r FlatRepositoryImpl) FindByID(id string) (model.Flat, error) {
	var flat model.Flat
	result := r.db.Where("id = ?", id).First(&flat)
	if result.Error != nil {
		return model.Flat{}, result.Error
	}
	return flat, nil
}

func (r FlatRepositoryImpl) FindAll() ([]model.Flat, error) {
	var flats []model.Flat
	result := r.db.Find(&flats)
	if result.Error != nil {
		return nil, result.Error
	}
	return flats, nil
}

func (r FlatRepositoryImpl) Create(flat model.Flat) (model.Flat, error) {
	if err := r.db.Create(&flat).Error; err != nil {
		return model.Flat{}, err
	}

	return flat, nil
}

func (r FlatRepositoryImpl) ExistsByID(id string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Flat{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
