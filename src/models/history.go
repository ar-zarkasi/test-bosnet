package models

import (
	"time"
)

type History struct {
    SzTransactionId string    `gorm:"type:nvarchar(50);primaryKey;column:szTransactionId" json:"szTransactionId"`
    SzAccountId     string    `gorm:"type:nvarchar(50);primaryKey;column:szAccountId" json:"szAccountId"`
    SzCurrencyId    string    `gorm:"type:nvarchar(50);primaryKey;column:szCurrencyId" json:"szCurrencyId"`
    DtmTransaction  time.Time `gorm:"type:datetime;not null;column:dtmTransaction" json:"dtmTransaction"`
    DecAmount       float64   `gorm:"type:decimal(30,8);not null;column:decAmount" json:"decAmount"`
    SzNote          string    `gorm:"type:nvarchar(255);not null;column:szNote" json:"szNote"`
}

func (t *History) TableName() string {
    return "BOS_History"
}
