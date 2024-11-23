package models

type Balance struct {
	SzAccountId string `gorm:"type:nvarchar(50);primaryKey;column:szAccountId" json:"szAccountId"`
	SzCurrencyId string `gorm:"type:nvarchar(50);primaryKey;column:szCurrencyId" json:"szCurrencyId"`
	DecAmount *float64 `gorm:"type:decimal(30,8);default:18;column:decAmount" json:"decAmount"`
}

func (t *Balance) TableName() string {
	return "BOS_Balance"
}