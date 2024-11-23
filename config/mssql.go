package config

import (
	"app/src/models"
	"app/utils"
	"fmt"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func runSqlServer() {
	globalDb = connectSqlServer()
	MigrateEntityModels(globalDb, &models.Config{}, &models.Counter{}, &models.Balance{}, &models.History{})
	// stringPathModel := "app/src/models/"
	// MigrateModels(globalDb, stringPathModel)
}

func connectSqlServer() *gorm.DB {
	// connect to PostgreSql
	var (
		host = os.Getenv("DB_HOST")
		user = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname = os.Getenv("DB_NAME")
		port = os.Getenv("DB_PORT")
		timezone = os.Getenv("DB_TIMEZONE")
	)

	// set default timezone if not set
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable", user, password, host, port, dbname)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	utils.ErrorFatal(err)

	return db
}

func getDBSQLServer() *gorm.DB {
	return globalDb
}