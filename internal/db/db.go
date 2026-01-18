package db

import (
	"fmt"
	"indexer/internal/config"
	"indexer/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	err = db.AutoMigrate(
		&model.BTCBlock{}, &model.BTCTransaction{},
		&model.ETHBlock{}, &model.ETHTransaction{},
		&model.IndexerState{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
