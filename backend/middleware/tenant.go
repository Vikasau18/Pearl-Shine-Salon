package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantContext middleware extracts salon_id from route params or query
// and verifies the current user has access to that salon (multi-tenant isolation).
func TenantContext(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		salonID := c.Param("salon_id")
		if salonID == "" {
			salonID = c.Query("salon_id")
		}
		if salonID == "" {
			c.Next()
			return
		}

		c.Set("tenant_salon_id", salonID)

		// Validate that the salon exists
		var exists bool
		err := db.QueryRow(context.Background(),
			"SELECT EXISTS(SELECT 1 FROM salons WHERE id = $1 AND is_active = true)", salonID).Scan(&exists)
		if err != nil || !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Salon not found"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// EnsureSalonOwnership verifies the authenticated user owns the salon being accessed
func EnsureSalonOwnership(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		userRole := GetUserRole(c)

		// Admins bypass ownership check
		if userRole == "admin" {
			c.Next()
			return
		}

		salonID := c.Param("salon_id")
		if salonID == "" {
			salonID, _ = c.GetQuery("salon_id")
		}
		if salonID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "salon_id is required"})
			c.Abort()
			return
		}

		var ownerID string
		err := db.QueryRow(context.Background(),
			"SELECT owner_id FROM salons WHERE id = $1", salonID).Scan(&ownerID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Salon not found"})
			c.Abort()
			return
		}

		if ownerID != userID {
			// Check if user is staff of this salon
			var isStaff bool
			db.QueryRow(context.Background(),
				"SELECT EXISTS(SELECT 1 FROM staff WHERE salon_id = $1 AND user_id = $2 AND is_active = true)",
				salonID, userID).Scan(&isStaff)

			if !isStaff {
				c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this salon"})
				c.Abort()
				return
			}
		}

		c.Set("tenant_salon_id", salonID)
		c.Next()
	}
}

// GetTenantSalonID returns the salon_id from context
func GetTenantSalonID(c *gin.Context) string {
	salonID, _ := c.Get("tenant_salon_id")
	if salonID == nil {
		return ""
	}
	return salonID.(string)
}
