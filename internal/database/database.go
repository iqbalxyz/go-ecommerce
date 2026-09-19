package database

import (
	"database/sql"
	"go-ecommerce/internal/config"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
)

var DB *sql.DB

func Connect() (*gorm.DB, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		DryRun: false,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns())
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns())
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime())
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime())

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	DB = sqlDB

	return db, nil
}
