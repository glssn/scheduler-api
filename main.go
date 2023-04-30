package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api"
	"github.com/glssn/scheduler-api/config"
	"github.com/glssn/scheduler-api/initializers"
)

func init() {
	config.LoadEnvVariables()
	initializers.Logger()
	initializers.ConnectToDB()
	initializers.MigrateDatabase()
	initializers.PopulateBankHolidays()
	go initializers.SyncBankHolidays()
}

func main() {

	app := gin.New()
	app.Use(gin.Recovery())
	config := cors.DefaultConfig()

	// Set the AllowOrigins field to a list of domains from the ALLOWED_ORIGINS environment variable
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		config.AllowOrigins = strings.Split(allowedOrigins, ",")
	} else {
		// If ALLOWED_ORIGINS is not set or is set to an empty string, allow requests from all origins
		config.AllowOrigins = []string{"*"}
	}

	config.AllowCredentials = true
	app.Use(cors.New(config))

	api.Routes(app)

	log.Fatal(app.Run("localhost:3000"))
}
