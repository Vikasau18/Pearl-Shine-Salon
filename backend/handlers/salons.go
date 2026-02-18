package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SalonHandler struct {
	DB *pgxpool.Pool
}

func NewSalonHandler(db *pgxpool.Pool) *SalonHandler {
	return &SalonHandler{DB: db}
}

func (h *SalonHandler) ListSalons(c *gin.Context) {
	city := c.Query("city")
	search := c.Query("search")
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "50") // km
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit
	userID := middleware.GetUserID(c)

	var args []interface{}
	argIdx := 1

	// Base query
	selectCols := "s.id, s.owner_id, s.name, s.address, s.city, COALESCE(s.state,''), COALESCE(s.zip_code,''), COALESCE(s.lat,0), COALESCE(s.lng,0), COALESCE(s.phone,''), COALESCE(s.email,''), COALESCE(s.description,''), s.rating, s.total_reviews, COALESCE(s.image_url,''), s.is_active, s.opening_time::text, s.closing_time::text, s.created_at, s.updated_at"

	// Distance calculation if lat/lng provided
	distanceCol := ""
	if latStr != "" && lngStr != "" {
		lat, _ := strconv.ParseFloat(latStr, 64)
		lng, _ := strconv.ParseFloat(lngStr, 64)
		radius, _ := strconv.ParseFloat(radiusStr, 64)
		_ = radius

		distanceCol = fmt.Sprintf(", (6371 * acos(cos(radians($%d)) * cos(radians(s.lat)) * cos(radians(s.lng) - radians($%d)) + sin(radians($%d)) * sin(radians(s.lat)))) AS distance", argIdx, argIdx+1, argIdx+2)
		args = append(args, lat, lng, lat)
		argIdx += 3
	}

	// Favorite check
	favCol := ""
	if userID != "" {
		favCol = fmt.Sprintf(", EXISTS(SELECT 1 FROM favorites WHERE user_id = $%d AND salon_id = s.id) AS is_favorited", argIdx)
		args = append(args, userID)
		argIdx++
	}

	query := fmt.Sprintf("SELECT %s%s%s FROM salons s WHERE s.is_active = true", selectCols, distanceCol, favCol)

	// Filters
	if city != "" {
		query += fmt.Sprintf(" AND LOWER(s.city) = LOWER($%d)", argIdx)
		args = append(args, city)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(s.name) LIKE LOWER($%d) OR LOWER(s.description) LIKE LOWER($%d))", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Distance filter
	if latStr != "" && lngStr != "" {
		radius, _ := strconv.ParseFloat(radiusStr, 64)
		query += fmt.Sprintf(" AND (6371 * acos(cos(radians($%d)) * cos(radians(s.lat)) * cos(radians(s.lng) - radians($%d)) + sin(radians($%d)) * sin(radians(s.lat)))) < $%d", argIdx, argIdx+1, argIdx+2, argIdx+3)
		args = append(args, args[0], args[1], args[0], radius)
		argIdx += 4
	}

	// Order
	if distanceCol != "" {
		query += " ORDER BY distance ASC"
	} else {
		query += " ORDER BY s.rating DESC, s.total_reviews DESC"
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM salons s WHERE s.is_active = true")
	var totalCount int
	h.DB.QueryRow(context.Background(), countQuery).Scan(&totalCount)

	// Pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := h.DB.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch salons", "details": err.Error()})
		return
	}
	defer rows.Close()

	var salons []models.Salon
	for rows.Next() {
		var s models.Salon
		scanArgs := []interface{}{
			&s.ID, &s.OwnerID, &s.Name, &s.Address, &s.City, &s.State, &s.ZipCode,
			&s.Lat, &s.Lng, &s.Phone, &s.Email, &s.Description, &s.Rating,
			&s.TotalReviews, &s.ImageURL, &s.IsActive, &s.OpeningTime, &s.ClosingTime,
			&s.CreatedAt, &s.UpdatedAt,
		}
		if distanceCol != "" {
			var dist float64
			scanArgs = append(scanArgs, &dist)
			s.Distance = &dist
		}
		if userID != "" {
			var fav bool
			scanArgs = append(scanArgs, &fav)
			s.IsFavorited = &fav
		}
		if err := rows.Scan(scanArgs...); err != nil {
			continue
		}
		salons = append(salons, s)
	}

	if salons == nil {
		salons = []models.Salon{}
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       salons,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: int(math.Ceil(float64(totalCount) / float64(limit))),
	})
}

func (h *SalonHandler) GetSalon(c *gin.Context) {
	salonID := c.Param("id")
	userID := middleware.GetUserID(c)

	var s models.Salon
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, owner_id, name, address, city, COALESCE(state,''), COALESCE(zip_code,''),
		 COALESCE(lat,0), COALESCE(lng,0), COALESCE(phone,''), COALESCE(email,''), 
		 COALESCE(description,''), rating, total_reviews, COALESCE(image_url,''), 
		 is_active, opening_time::text, closing_time::text, created_at, updated_at 
		 FROM salons WHERE id = $1`, salonID,
	).Scan(&s.ID, &s.OwnerID, &s.Name, &s.Address, &s.City, &s.State, &s.ZipCode,
		&s.Lat, &s.Lng, &s.Phone, &s.Email, &s.Description, &s.Rating,
		&s.TotalReviews, &s.ImageURL, &s.IsActive, &s.OpeningTime, &s.ClosingTime,
		&s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Salon not found"})
		return
	}

	if userID != "" {
		var fav bool
		h.DB.QueryRow(context.Background(),
			"SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND salon_id = $2)",
			userID, salonID).Scan(&fav)
		s.IsFavorited = &fav
	}

	c.JSON(http.StatusOK, s)
}

func (h *SalonHandler) CreateSalon(c *gin.Context) {
	var req models.CreateSalonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerID := middleware.GetUserID(c)
	if req.OpeningTime == "" {
		req.OpeningTime = "09:00"
	}
	if req.ClosingTime == "" {
		req.ClosingTime = "21:00"
	}

	var s models.Salon
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO salons (owner_id, name, address, city, state, zip_code, lat, lng, phone, email, description, image_url, opening_time, closing_time) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		 RETURNING id, owner_id, name, address, city, COALESCE(state,''), COALESCE(zip_code,''),
		 COALESCE(lat,0), COALESCE(lng,0), COALESCE(phone,''), COALESCE(email,''), 
		 COALESCE(description,''), rating, total_reviews, COALESCE(image_url,''), 
		 is_active, opening_time::text, closing_time::text, created_at, updated_at`,
		ownerID, req.Name, req.Address, req.City, req.State, req.ZipCode,
		req.Lat, req.Lng, req.Phone, req.Email, req.Description, req.ImageURL,
		req.OpeningTime, req.ClosingTime,
	).Scan(&s.ID, &s.OwnerID, &s.Name, &s.Address, &s.City, &s.State, &s.ZipCode,
		&s.Lat, &s.Lng, &s.Phone, &s.Email, &s.Description, &s.Rating,
		&s.TotalReviews, &s.ImageURL, &s.IsActive, &s.OpeningTime, &s.ClosingTime,
		&s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create salon", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func (h *SalonHandler) UpdateSalon(c *gin.Context) {
	salonID := c.Param("salon_id")
	var req models.CreateSalonRequest
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

	var s models.Salon
	err := h.DB.QueryRow(context.Background(),
		`UPDATE salons SET name=$1, address=$2, city=$3, state=$4, zip_code=$5, lat=$6, lng=$7,
		 phone=$8, email=$9, description=$10, image_url=$11, 
		 opening_time=$12, closing_time=$13, updated_at=NOW()
		 WHERE id=$14
		 RETURNING id, owner_id, name, address, city, COALESCE(state,''), COALESCE(zip_code,''),
		 COALESCE(lat,0), COALESCE(lng,0), COALESCE(phone,''), COALESCE(email,''), 
		 COALESCE(description,''), rating, total_reviews, COALESCE(image_url,''), 
		 is_active, opening_time::text, closing_time::text, created_at, updated_at`,
		req.Name, req.Address, req.City, req.State, req.ZipCode, req.Lat, req.Lng,
		req.Phone, req.Email, req.Description, req.ImageURL,
		req.OpeningTime, req.ClosingTime, salonID,
	).Scan(&s.ID, &s.OwnerID, &s.Name, &s.Address, &s.City, &s.State, &s.ZipCode,
		&s.Lat, &s.Lng, &s.Phone, &s.Email, &s.Description, &s.Rating,
		&s.TotalReviews, &s.ImageURL, &s.IsActive, &s.OpeningTime, &s.ClosingTime,
		&s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update salon"})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *SalonHandler) GetGallery(c *gin.Context) {
	salonID := c.Param("id")

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, salon_id, image_url, COALESCE(caption,''), sort_order, created_at 
		 FROM salon_gallery WHERE salon_id = $1 ORDER BY sort_order`, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gallery"})
		return
	}
	defer rows.Close()

	var gallery []models.SalonGallery
	for rows.Next() {
		var g models.SalonGallery
		rows.Scan(&g.ID, &g.SalonID, &g.ImageURL, &g.Caption, &g.SortOrder, &g.CreatedAt)
		gallery = append(gallery, g)
	}
	if gallery == nil {
		gallery = []models.SalonGallery{}
	}

	c.JSON(http.StatusOK, gallery)
}

// GetMySalons returns salons owned by the authenticated user (multi-tenant)
func (h *SalonHandler) GetMySalons(c *gin.Context) {
	ownerID := middleware.GetUserID(c)

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, owner_id, name, address, city, COALESCE(state,''), COALESCE(zip_code,''),
		 COALESCE(lat,0), COALESCE(lng,0), COALESCE(phone,''), COALESCE(email,''), 
		 COALESCE(description,''), rating, total_reviews, COALESCE(image_url,''), 
		 is_active, opening_time::text, closing_time::text, created_at, updated_at 
		 FROM salons WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch salons"})
		return
	}
	defer rows.Close()

	var salons []models.Salon
	for rows.Next() {
		var s models.Salon
		rows.Scan(&s.ID, &s.OwnerID, &s.Name, &s.Address, &s.City, &s.State, &s.ZipCode,
			&s.Lat, &s.Lng, &s.Phone, &s.Email, &s.Description, &s.Rating,
			&s.TotalReviews, &s.ImageURL, &s.IsActive, &s.OpeningTime, &s.ClosingTime,
			&s.CreatedAt, &s.UpdatedAt)
		salons = append(salons, s)
	}
	if salons == nil {
		salons = []models.Salon{}
	}

	c.JSON(http.StatusOK, salons)
}
