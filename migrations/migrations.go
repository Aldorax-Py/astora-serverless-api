package main

import (
	"astora/inits"
	"astora/models"
)

func init() {
	// LoadEnv()
	inits.LoadEnv()
	inits.DBInit()
}

func main() {
	// Migrate User and Services tables first
	inits.DB.AutoMigrate(&models.User{})
	inits.DB.AutoMigrate(&models.Services{})

	// Migrate ServiceLogs and ApiResponses tables
	inits.DB.AutoMigrate(&models.ServiceLogs{})
	inits.DB.AutoMigrate(&models.ApiResponses{})
}
