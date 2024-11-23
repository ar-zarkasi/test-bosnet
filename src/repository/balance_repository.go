package repository

import (
	"app/src/http/request"
	interfaces "app/src/interface"
	"app/src/models"

	"gorm.io/gorm"
)

type BalanceRepository struct {
	Db *gorm.DB
}

func NewBalance(db *gorm.DB) interfaces.BalanceInterface {
	return &BalanceRepository{Db: db}
}

func (r *BalanceRepository) Create(data request.DataBalance, model *models.Balance) error {
	model.SzAccountId = data.Account
	model.DecAmount = &data.Balance
	model.SzCurrencyId = data.Currency
	return r.Db.Create(&model).Error
}
func (r *BalanceRepository) UpdateSaldo(balance float64, model *models.Balance) error {
	model.DecAmount = &balance
	return r.Db.Save(&model).Error
}
func (r *BalanceRepository) FindOne(account string, currency string) (*models.Balance, error) {
	var model models.Balance
	err := r.Db.First(&model,"szAccountId = ? AND szCurrencyId = ?", account, currency).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *BalanceRepository) FindAllByAccount(account string) ([]models.Balance, error) {
	models := make([]models.Balance,0)
	err := r.Db.Where("szAccountId = ?", account).Find(&models).Error
	if err != nil {
		return nil, err
	}

	return models, nil
}
