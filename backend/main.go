package main

import (
	"context"
	"log"

	"saloon-backend/config"
	"saloon-backend/database"
	"saloon-backend/middleware"
	"saloon-backend/routes"
	"saloon-backend/services"

	webpush "github.com/SherClockHolmes/webpush-go"
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

	// Auto-generate VAPID keys if not set
	if cfg.VAPIDPublicKey == "" || cfg.VAPIDPrivateKey == "" {
		log.Println("⚠️  VAPID keys not found, generating new keys...")
		privKey, pubKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			log.Fatalf("❌ Failed to generate VAPID keys: %v", err)
		}
		cfg.VAPIDPublicKey = pubKey
		cfg.VAPIDPrivateKey = privKey
		log.Printf("🔑 VAPID Public Key:  %s", pubKey)
		log.Printf("🔑 VAPID Private Key: %s", privKey)
		log.Println("💡 Set VAPID_PUBLIC_KEY and VAPID_PRIVATE_KEY env vars to persist these keys.")
	}

	// Initialize push service
	pushService := services.NewPushService(db, cfg.VAPIDPublicKey, cfg.VAPIDPrivateKey, cfg.VAPIDContact)

	// Start background reminder scheduler
	scheduler := services.NewScheduler(db, pushService)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduler.Start(ctx)

	// Create Gin router
	r := gin.Default()

	// Apply global middleware
	r.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	// Initialize Cloudinary service
	cloudinaryService, err := services.NewCloudinaryService(cfg.CloudinaryURL)
	if err != nil {
		log.Printf("⚠️  Failed to initialize Cloudinary service: %v", err)
	}

	// Register routes
	routes.RegisterRoutes(r, db, pushService, scheduler, cloudinaryService)

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
