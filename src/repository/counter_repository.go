package repository

import (
	"app/src/http/request"
	interfaces "app/src/interface"
	"app/src/models"

	"gorm.io/gorm"
)


type CounterRepository struct {
	Db *gorm.DB
}

func NewCounter(db *gorm.DB) (interfaces.CounterInterface) {
	return &CounterRepository{Db: db}
}

func (r *CounterRepository) Create(data request.DataCounter, model *models.Counter) error {
	model.SzCounterId = data.CounterId
	model.ILastNumber = data.LastNumber
	return r.Db.Create(&model).Error
}

func (r *CounterRepository) Update(model *models.Counter) error {
	return r.Db.Save(&model).Error
}

func (r *CounterRepository) Find(filter *map[string]interface{}) []models.Counter {
	lists := make([]models.Counter, 0)
	tempdb := r.Db
	for key, value := range *filter {
		tempdb = tempdb.Where(key, value)
	}
	err := tempdb.Find(&lists).Error
	if err != nil {
		return nil
	}

	return lists
}