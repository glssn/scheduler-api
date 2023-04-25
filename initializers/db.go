package initializers

import (
	"log"
	"os"

	"github.com/glssn/scheduler-api/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	dbURL := os.Getenv("DATABASE_URI")
	if dbURL == "" {
		log.Fatal("DATABASE_URI environment variable unset")
		os.Exit(2)
	}
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}
	log.Println("Connected Successfully to Database")
}

func MigrateDatabase() {
	// get a logger instance
	logger := Logger()

	// log the start of the function
	logger.Println("MigrateDatabase: start")

	DB.AutoMigrate(
		models.Event{},
		models.EventMeta{},
		models.User{})

	// log the end of the function
	logger.Println("MigrateDatabase: end")
}
