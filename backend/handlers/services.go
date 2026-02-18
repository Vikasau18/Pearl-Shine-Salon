package handlers

import (
	"context"
	"net/http"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceHandler struct {
	DB *pgxpool.Pool
}

func NewServiceHandler(db *pgxpool.Pool) *ServiceHandler {
	return &ServiceHandler{DB: db}
}

func (h *ServiceHandler) ListServices(c *gin.Context) {
	salonID := c.Param("id")
	if salonID == "" {
		salonID = c.Param("salon_id")
	}
	category := c.Query("category")

	query := `SELECT id, salon_id, name, COALESCE(description,''), price, duration_minutes, 
			  buffer_minutes, COALESCE(category,''), is_active, created_at, updated_at 
			  FROM services WHERE salon_id = $1 AND is_active = true`
	args := []interface{}{salonID}

	if category != "" {
		query += " AND LOWER(category) = LOWER($2)"
		args = append(args, category)
	}
	query += " ORDER BY category, name"

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		rows.Scan(&s.ID, &s.SalonID, &s.Name, &s.Description, &s.Price, &s.DurationMinutes,
			&s.BufferMinutes, &s.Category, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		services = append(services, s)
	}
	if services == nil {
		services = []models.Service{}
	}

	c.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	salonID := c.Param("salon_id")
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	ownerID := middleware.GetUserID(c)
	var existingOwner string
	h.DB.QueryRow(context.Background(), "SELECT owner_id FROM salons WHERE id = $1", salonID).Scan(&existingOwner)
	if existingOwner != ownerID && middleware.GetUserRole(c) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this salon"})
		return
	}

	var s models.Service
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO services (salon_id, name, description, price, duration_minutes, buffer_minutes, category)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, salon_id, name, COALESCE(description,''), price, duration_minutes, buffer_minutes, COALESCE(category,''), is_active, created_at, updated_at`,
		salonID, req.Name, req.Description, req.Price, req.DurationMinutes, req.BufferMinutes, req.Category,
	).Scan(&s.ID, &s.SalonID, &s.Name, &s.Description, &s.Price, &s.DurationMinutes,
		&s.BufferMinutes, &s.Category, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	serviceID := c.Param("service_id")
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s models.Service
	err := h.DB.QueryRow(context.Background(),
		`UPDATE services SET name=$1, description=$2, price=$3, duration_minutes=$4, buffer_minutes=$5, category=$6, updated_at=NOW()
		 WHERE id=$7
		 RETURNING id, salon_id, name, COALESCE(description,''), price, duration_minutes, buffer_minutes, COALESCE(category,''), is_active, created_at, updated_at`,
		req.Name, req.Description, req.Price, req.DurationMinutes, req.BufferMinutes, req.Category, serviceID,
	).Scan(&s.ID, &s.SalonID, &s.Name, &s.Description, &s.Price, &s.DurationMinutes,
		&s.BufferMinutes, &s.Category, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	serviceID := c.Param("service_id")

	_, err := h.DB.Exec(context.Background(),
		"UPDATE services SET is_active = false, updated_at = NOW() WHERE id = $1", serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
