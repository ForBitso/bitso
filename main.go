package main

import (
	"log"

	"go-shop/config"
	"go-shop/database"
	"go-shop/routes"
)

func main() {

	// Load configuration
	cfg := config.Load()

	// Connect to database
	database.ConnectDB(cfg)
	database.Migrate()

	// Connect to Redis
	database.ConnectRedis(cfg)

	// Setup routes
	router := routes.SetupRoutes(cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

//test to git
