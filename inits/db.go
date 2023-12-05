package inits

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInit() {
	dsn := os.Getenv("RAILWAY_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Set up connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to set up connection pooling")
	}

	// Set the maximum number of idle connections in the pool
	sqlDB.SetMaxIdleConns(10)

	// Set the maximum number of open connections (both idle and in-use) in the pool
	sqlDB.SetMaxOpenConns(100)

	DB = db
}
