package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	env := os.Getenv("SCHEDULER_API_ENV")
	// load environment variables from dotenv
	if env == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}
