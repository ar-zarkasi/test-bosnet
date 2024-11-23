package repository

import (
	interfaces "app/src/interface"
	"app/src/models"
	"fmt"

	"gorm.io/gorm"
)


type SettingRepository struct {
	Db *gorm.DB
}

func NewSettings(db *gorm.DB) interfaces.SettingInterface {
	return &SettingRepository{Db: db}
}

func (r *SettingRepository) Create(key string, value string) (*models.Config, error) {
	model := models.Config{
		ConfigKey: key,
		ConfigValue: value,
	}
	err := r.Db.Create(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}
func (r *SettingRepository) Update(value string, model *models.Config) error {
	model.ConfigValue = value
	err := r.Db.Save(&model).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *SettingRepository) Find(Key string) *models.Config {
	var model models.Config
	err := r.Db.First(&model, "config_key = ?",Key).Error
	if err != nil {
		fmt.Println("Setting Not Found")
		return nil
	}
	return &model
}

func (r *SettingRepository) BeginTransaction() *gorm.DB {
	return r.Db.Begin()
}