package models

type Config struct {
	ID uint `gorm:"bigint;primaryKey;autoIncrement" json:"id"`
	ConfigKey string `gorm:"nvarchar(50);not null;index;" json:"key"`
	ConfigValue string `gorm:"nvarchar(255);nullable" json:"value"`
}

func (t *Config) TableName() string {
	return "settings"
}