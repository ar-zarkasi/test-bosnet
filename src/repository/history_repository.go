package repository

import (
	interfaces "app/src/interface"
	"app/src/models"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HistoryRepository struct {
	Db *gorm.DB
}

func NewHistory(db *gorm.DB) interfaces.HistoryInterface {
	return &HistoryRepository{Db: db}
}

func (r *HistoryRepository) Create(model *models.History) error {
	return r.Db.Create(&model).Error
}
func (r *HistoryRepository) FindTransaction(filter *map[string]interface{}, sortBy string, Descending bool) ([]models.History, error) {
	lists := make([]models.History, 0)
	tempdb := r.Db
	for key, value := range *filter {
		if key == "dateBetween" {
			ds := value.(map[string]string)
			tempdb = tempdb.Where("dtmTransaction BETWEEN ? AND ?", ds["from"], ds["to"])
			continue
		}
		tempdb = tempdb.Where(key, value)
	}
	tempdb = tempdb.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: Descending})
	err := tempdb.Find(&lists).Error
	if err != nil {
		return nil, err
	}

	return lists, nil
}
func (r *HistoryRepository) BeginTransaction() *gorm.DB {
	return r.Db.Begin()
}
func (r *HistoryRepository) CountTransaction(transactionId string) int {
	count := int64(0)
	db := r.Db.Model(&models.History{}).Where("szTransactionId LIKE ?","%"+transactionId)
	err := db.Count(&count).Error
	if err != nil {
		fmt.Println("Counting Transaction Error", err)
		return 0
	}
	integerCount := int(count)
	return integerCount
}