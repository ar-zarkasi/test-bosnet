package models

type Counter struct {
	SzCounterId string `gorm:"type:nvarchar(50);primaryKey;column:szCounterId" json:"szCounterId"`
	ILastNumber uint `gorm:"type:bigint;not null;column:iLastNumber" json:"iLastNumber"`
}

func (t *Counter) TableName() string {
	return "BOS_Counter"
}
