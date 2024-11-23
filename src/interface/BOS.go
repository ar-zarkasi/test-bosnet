package interfaces

import (
	"app/src/http/request"
	"app/src/models"

	"gorm.io/gorm"
)

type SettingInterface interface {
	Create(key string, value string) (*models.Config, error)
	Update(value string, model *models.Config) error
	Find(Key string) *models.Config
	BeginTransaction() *gorm.DB
}

type CounterInterface interface {
	Create(data request.DataCounter, model *models.Counter) error
	Find(filter *map[string]interface{}) []models.Counter
	Update(model *models.Counter) error
}

type BalanceInterface interface {
	Create(data request.DataBalance, model *models.Balance) error
	UpdateSaldo(balance float64, model *models.Balance) error
	FindOne(account string, currency string) (*models.Balance, error)
	FindAllByAccount(account string) ([]models.Balance, error)
}

type HistoryInterface interface {
	Create(model *models.History) error
	FindTransaction(filter *map[string]interface{}, sortBy string, Descending bool) ([]models.History, error)
	BeginTransaction() *gorm.DB
	CountTransaction(transactionId string) int
}
