package handlers

import (
	"context"
	"fmt"
	"net/http"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffHandler struct {
	DB *pgxpool.Pool
}

func NewStaffHandler(db *pgxpool.Pool) *StaffHandler {
	return &StaffHandler{DB: db}
}

func (h *StaffHandler) ListStaff(c *gin.Context) {
	salonID := c.Param("id")
	if salonID == "" {
		salonID = c.Param("salon_id")
	}

	rows, err := h.DB.Query(context.Background(),
		`SELECT s.id, s.salon_id, s.user_id, s.name, s.role, COALESCE(s.specialization,''), 
		 COALESCE(s.bio,''), COALESCE(s.avatar_url,''), s.is_active, s.created_at, s.updated_at
		 FROM staff s WHERE s.salon_id = $1 AND s.is_active = true
		 ORDER BY s.name`, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff"})
		return
	}
	defer rows.Close()

	var staffList []models.Staff
	for rows.Next() {
		var s models.Staff
		rows.Scan(&s.ID, &s.SalonID, &s.UserID, &s.Name, &s.Role, &s.Specialization,
			&s.Bio, &s.AvatarURL, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)

		// Fetch services for this staff
		svcRows, _ := h.DB.Query(context.Background(),
			`SELECT sv.id, sv.salon_id, sv.name, COALESCE(sv.description,''), sv.price, sv.duration_minutes, 
			 sv.buffer_minutes, COALESCE(sv.category,''), sv.is_active, sv.created_at, sv.updated_at
			 FROM services sv 
			 JOIN staff_services ss ON ss.service_id = sv.id 
			 WHERE ss.staff_id = $1 AND sv.is_active = true`, s.ID)
		if svcRows != nil {
			for svcRows.Next() {
				var svc models.Service
				svcRows.Scan(&svc.ID, &svc.SalonID, &svc.Name, &svc.Description, &svc.Price,
					&svc.DurationMinutes, &svc.BufferMinutes, &svc.Category, &svc.IsActive,
					&svc.CreatedAt, &svc.UpdatedAt)
				s.Services = append(s.Services, svc)
			}
			svcRows.Close()
		}
		if s.Services == nil {
			s.Services = []models.Service{}
		}

		staffList = append(staffList, s)
	}
	if staffList == nil {
		staffList = []models.Staff{}
	}

	c.JSON(http.StatusOK, staffList)
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	salonID := c.Param("salon_id")
	var req models.CreateStaffRequest
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

	if req.Role == "" {
		req.Role = "stylist"
	}

	var s models.Staff
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO staff (salon_id, name, role, specialization, bio, avatar_url) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id, salon_id, user_id, name, role, COALESCE(specialization,''), COALESCE(bio,''), COALESCE(avatar_url,''), is_active, created_at, updated_at`,
		salonID, req.Name, req.Role, req.Specialization, req.Bio, req.AvatarURL,
	).Scan(&s.ID, &s.SalonID, &s.UserID, &s.Name, &s.Role, &s.Specialization,
		&s.Bio, &s.AvatarURL, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff", "details": err.Error()})
		return
	}

	// Assign services
	for _, svcID := range req.ServiceIDs {
		h.DB.Exec(context.Background(),
			"INSERT INTO staff_services (staff_id, service_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			s.ID, svcID)
	}

	// Assign default working hours (Mon-Fri 09:00-17:00, Sat 10:00-15:00, Sun Off)
	defaultHours := []struct {
		day   int
		start string
		end   string
		off   bool
	}{
		{0, "09:00", "17:00", true},  // Sunday (Off)
		{1, "09:00", "17:00", false}, // Monday
		{2, "09:00", "17:00", false}, // Tuesday
		{3, "09:00", "17:00", false}, // Wednesday
		{4, "09:00", "17:00", false}, // Thursday
		{5, "09:00", "17:00", false}, // Friday
		{6, "10:00", "15:00", false}, // Saturday
	}

	for _, dh := range defaultHours {
		_, err := h.DB.Exec(context.Background(),
			`INSERT INTO staff_working_hours (staff_id, day_of_week, start_time, end_time, is_off)
			 VALUES ($1, $2, $3, $4, $5)`,
			s.ID, dh.day, dh.start, dh.end, dh.off)
		if err != nil {
			// Log error but continue
			fmt.Printf("Failed to set default hours for staff %s day %d: %v\n", s.ID, dh.day, err)
		}
	}

	s.Services = []models.Service{}
	c.JSON(http.StatusCreated, s)
}

func (h *StaffHandler) UpdateStaff(c *gin.Context) {
	staffID := c.Param("staff_id")
	var req models.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var s models.Staff
	err := h.DB.QueryRow(context.Background(),
		`UPDATE staff SET name=$1, role=$2, specialization=$3, bio=$4, avatar_url=$5, updated_at=NOW()
		 WHERE id=$6
		 RETURNING id, salon_id, user_id, name, role, COALESCE(specialization,''), COALESCE(bio,''), COALESCE(avatar_url,''), is_active, created_at, updated_at`,
		req.Name, req.Role, req.Specialization, req.Bio, req.AvatarURL, staffID,
	).Scan(&s.ID, &s.SalonID, &s.UserID, &s.Name, &s.Role, &s.Specialization,
		&s.Bio, &s.AvatarURL, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff"})
		return
	}

	// Update service assignments
	if req.ServiceIDs != nil {
		h.DB.Exec(context.Background(), "DELETE FROM staff_services WHERE staff_id = $1", staffID)
		for _, svcID := range req.ServiceIDs {
			h.DB.Exec(context.Background(),
				"INSERT INTO staff_services (staff_id, service_id) VALUES ($1, $2)", staffID, svcID)
		}
	}

	s.Services = []models.Service{}
	c.JSON(http.StatusOK, s)
}

func (h *StaffHandler) DeleteStaff(c *gin.Context) {
	staffID := c.Param("staff_id")

	_, err := h.DB.Exec(context.Background(),
		"UPDATE staff SET is_active = false, updated_at = NOW() WHERE id = $1", staffID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff member removed"})
}

func (h *StaffHandler) GetWorkingHours(c *gin.Context) {
	staffID := c.Param("staff_id")

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, staff_id, day_of_week, start_time::text, end_time::text, is_off 
		 FROM staff_working_hours WHERE staff_id = $1 ORDER BY day_of_week`, staffID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch working hours"})
		return
	}
	defer rows.Close()

	var hours []models.StaffWorkingHours
	for rows.Next() {
		var h models.StaffWorkingHours
		rows.Scan(&h.ID, &h.StaffID, &h.DayOfWeek, &h.StartTime, &h.EndTime, &h.IsOff)
		hours = append(hours, h)
	}
	if hours == nil {
		hours = []models.StaffWorkingHours{}
	}

	c.JSON(http.StatusOK, hours)
}

func (h *StaffHandler) UpdateWorkingHours(c *gin.Context) {
	staffID := c.Param("staff_id")
	var req models.UpdateStaffWorkingHoursRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete existing and re-insert
	h.DB.Exec(context.Background(), "DELETE FROM staff_working_hours WHERE staff_id = $1", staffID)

	for _, wh := range req.Hours {
		h.DB.Exec(context.Background(),
			`INSERT INTO staff_working_hours (staff_id, day_of_week, start_time, end_time, is_off)
			 VALUES ($1, $2, $3, $4, $5) ON CONFLICT (staff_id, day_of_week) DO UPDATE 
			 SET start_time = $3, end_time = $4, is_off = $5`,
			staffID, wh.DayOfWeek, wh.StartTime, wh.EndTime, wh.IsOff)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Working hours updated"})
}

// GetStaffForService returns staff members who can perform a specific service
func (h *StaffHandler) GetStaffForService(c *gin.Context) {
	salonID := c.Param("id")
	serviceID := c.Query("service_id")

	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_id query parameter required"})
		return
	}

	rows, err := h.DB.Query(context.Background(),
		`SELECT s.id, s.salon_id, s.user_id, s.name, s.role, COALESCE(s.specialization,''), 
		 COALESCE(s.bio,''), COALESCE(s.avatar_url,''), s.is_active, s.created_at, s.updated_at
		 FROM staff s
		 JOIN staff_services ss ON ss.staff_id = s.id
		 WHERE s.salon_id = $1 AND ss.service_id = $2 AND s.is_active = true
		 ORDER BY s.name`, salonID, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staff"})
		return
	}
	defer rows.Close()

	var staffList []models.Staff
	for rows.Next() {
		var s models.Staff
		rows.Scan(&s.ID, &s.SalonID, &s.UserID, &s.Name, &s.Role, &s.Specialization,
			&s.Bio, &s.AvatarURL, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		s.Services = []models.Service{}
		staffList = append(staffList, s)
	}
	if staffList == nil {
		staffList = []models.Staff{}
	}

	c.JSON(http.StatusOK, staffList)
}
