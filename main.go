package main

import (
	"products-api/config"
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
	r.Run(":8080")
}
