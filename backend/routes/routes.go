package routes

import (
	"saloon-backend/handlers"
	"saloon-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	salonHandler := handlers.NewSalonHandler(db)
	serviceHandler := handlers.NewServiceHandler(db)
	staffHandler := handlers.NewStaffHandler(db)
	appointmentHandler := handlers.NewAppointmentHandler(db)
	reviewHandler := handlers.NewReviewHandler(db)
	paymentHandler := handlers.NewPaymentHandler(db)
	favoriteHandler := handlers.NewFavoriteHandler(db)
	promoHandler := handlers.NewPromoHandler(db)
	waitlistHandler := handlers.NewWaitlistHandler(db)
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	notificationHandler := handlers.NewNotificationHandler(db)

	api := r.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "saloon-backend"})
	})

	// ─────────────────────────────────────────────
	// AUTH ROUTES
	// ─────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthRequired(), authHandler.GetMe)
		auth.PUT("/profile", middleware.AuthRequired(), authHandler.UpdateProfile)
	}

	// ─────────────────────────────────────────────
	// PUBLIC SALON ROUTES (with optional auth for favorites check)
	// ─────────────────────────────────────────────
	salons := api.Group("/salons")
	{
		salons.GET("", middleware.OptionalAuth(), salonHandler.ListSalons)
		salons.GET("/:id", middleware.OptionalAuth(), salonHandler.GetSalon)
		salons.GET("/:id/services", serviceHandler.ListServices)
		salons.GET("/:id/staff", staffHandler.ListStaff)
		salons.GET("/:id/staff-for-service", staffHandler.GetStaffForService)
		salons.GET("/:id/reviews", reviewHandler.GetSalonReviews)
		salons.GET("/:id/gallery", salonHandler.GetGallery)
	}

	// ─────────────────────────────────────────────
	// AUTHENTICATED CUSTOMER ROUTES
	// ─────────────────────────────────────────────
	customer := api.Group("")
	customer.Use(middleware.AuthRequired())
	{
		// Appointments
		customer.POST("/appointments", appointmentHandler.BookAppointment)
		customer.GET("/appointments", appointmentHandler.GetMyAppointments)
		customer.PUT("/appointments/:id/reschedule", appointmentHandler.RescheduleAppointment)
		customer.PUT("/appointments/:id/cancel", appointmentHandler.CancelAppointment)
		customer.GET("/appointments/available-slots", appointmentHandler.GetAvailableSlots)

		// Reviews
		customer.POST("/reviews", reviewHandler.CreateReview)

		// Favorites
		customer.POST("/favorites/:salon_id", favoriteHandler.ToggleFavorite)
		customer.GET("/favorites", favoriteHandler.GetMyFavorites)

		// Promo validation
		customer.GET("/promos/validate", promoHandler.ValidatePromo)

		// Waitlist
		customer.POST("/waitlist", waitlistHandler.JoinWaitlist)
		customer.DELETE("/waitlist/:id", waitlistHandler.LeaveWaitlist)

		// Notifications
		customer.GET("/notifications", notificationHandler.GetNotifications)
		customer.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
		customer.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		customer.PUT("/notifications/read-all", notificationHandler.MarkAllRead)
	}

	// ─────────────────────────────────────────────
	// SALON OWNER / DASHBOARD ROUTES (Multi-tenant)
	// ─────────────────────────────────────────────
	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthRequired())
	dashboard.Use(middleware.RoleRequired("salon_owner", "admin", "staff"))
	{
		// Salon management
		dashboard.POST("/salons", salonHandler.CreateSalon)
		dashboard.GET("/salons", salonHandler.GetMySalons)
		dashboard.PUT("/salons/:salon_id", salonHandler.UpdateSalon)

		// Per-salon management (tenant-scoped)
		salon := dashboard.Group("/salons/:salon_id")
		salon.Use(middleware.EnsureSalonOwnership(db))
		{
			// Services
			salon.GET("/services", serviceHandler.ListServices)
			salon.POST("/services", serviceHandler.CreateService)
			salon.PUT("/services/:service_id", serviceHandler.UpdateService)
			salon.DELETE("/services/:service_id", serviceHandler.DeleteService)

			// Staff
			salon.GET("/staff", staffHandler.ListStaff)
			salon.POST("/staff", staffHandler.CreateStaff)
			salon.PUT("/staff/:staff_id", staffHandler.UpdateStaff)
			salon.DELETE("/staff/:staff_id", staffHandler.DeleteStaff)
			salon.GET("/staff/:staff_id/hours", staffHandler.GetWorkingHours)
			salon.PUT("/staff/:staff_id/hours", staffHandler.UpdateWorkingHours)

			// Appointments
			salon.GET("/appointments", appointmentHandler.GetSalonAppointments)
			salon.PUT("/appointments/:id/no-show", appointmentHandler.MarkNoShow)
			salon.PUT("/appointments/:id/complete", appointmentHandler.CompleteAppointment)

			// Payments
			salon.GET("/payments", paymentHandler.GetSalonPayments)
			salon.POST("/payments", paymentHandler.ProcessPayment)
			salon.GET("/payments/:id/receipt", paymentHandler.GetReceipt)

			// Promo codes
			salon.POST("/promos", promoHandler.CreatePromo)
			salon.GET("/promos", promoHandler.GetSalonPromos)
			salon.DELETE("/promos/:promo_id", promoHandler.DeletePromo)

			// Waitlist
			salon.GET("/waitlist", waitlistHandler.GetSalonWaitlist)

			// Analytics
			salon.GET("/analytics", analyticsHandler.GetDashboard)
		}
	}
}
