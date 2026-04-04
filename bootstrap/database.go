package bootstrap

import (
	"log"

	"github.com/mhasnanr/e-wallet/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	dsn := GetEnv("CONNECTION_STRING", "")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database")
	}

	Log.Infow("connected to database...")

	database.AutoMigrate(&models.User{}, &models.UserSession{})

	return database, err
}
