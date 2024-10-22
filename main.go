package main

import (
	"os"
	"products-api/config"
	_ "products-api/docs"
	"products-api/global"
	"products-api/helpers"
	"products-api/routes"
)

func main() {
	// Initialize Logger
	helpers.InitLogger()

	// Initialize Configurations (DB, Redis)
	config.InitConfig()

	// Set up Routes
	r := routes.SetupRouter()

	// Start the server
	helpers.Info("Starting server on port 8080", nil)
	appPort := os.Getenv(global.PORT)
	if appPort != "" {
		r.Run(":%v", appPort)
	} else {
		r.Run(":8080")
	}
}
