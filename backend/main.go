package main

import (
	"log"

	"saloon-backend/config"
	"saloon-backend/database"
	"saloon-backend/middleware"
	"saloon-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set JWT secret
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create Gin router
	r := gin.Default()

	// Apply global middleware
	r.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	// Register routes
	routes.RegisterRoutes(r, db)

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
