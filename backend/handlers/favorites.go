package handlers

import (
	"context"
	"net/http"

	"saloon-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FavoriteHandler struct {
	DB *pgxpool.Pool
}

func NewFavoriteHandler(db *pgxpool.Pool) *FavoriteHandler {
	return &FavoriteHandler{DB: db}
}

func (h *FavoriteHandler) ToggleFavorite(c *gin.Context) {
	salonID := c.Param("salon_id")
	userID := middleware.GetUserID(c)

	// Check if already favorited
	var exists bool
	h.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND salon_id = $2)",
		userID, salonID).Scan(&exists)

	if exists {
		// Remove
		h.DB.Exec(context.Background(),
			"DELETE FROM favorites WHERE user_id = $1 AND salon_id = $2", userID, salonID)
		c.JSON(http.StatusOK, gin.H{"favorited": false, "message": "Removed from favorites"})
	} else {
		// Add
		h.DB.Exec(context.Background(),
			"INSERT INTO favorites (user_id, salon_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			userID, salonID)
		c.JSON(http.StatusOK, gin.H{"favorited": true, "message": "Added to favorites"})
	}
}

func (h *FavoriteHandler) GetMyFavorites(c *gin.Context) {
	userID := middleware.GetUserID(c)

	rows, err := h.DB.Query(context.Background(),
		`SELECT s.id, s.owner_id, s.name, s.address, s.city, COALESCE(s.state,''), 
		 COALESCE(s.zip_code,''), COALESCE(s.lat,0), COALESCE(s.lng,0), 
		 COALESCE(s.phone,''), COALESCE(s.email,''), COALESCE(s.description,''), 
		 s.rating, s.total_reviews, COALESCE(s.image_url,''), s.is_active,
		 s.opening_time::text, s.closing_time::text, s.created_at, s.updated_at
		 FROM salons s
		 JOIN favorites f ON f.salon_id = s.id
		 WHERE f.user_id = $1
		 ORDER BY f.created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}
	defer rows.Close()

	type FavoriteSalon struct {
		ID           string  `json:"id"`
		OwnerID      string  `json:"owner_id"`
		Name         string  `json:"name"`
		Address      string  `json:"address"`
		City         string  `json:"city"`
		State        string  `json:"state"`
		ZipCode      string  `json:"zip_code"`
		Lat          float64 `json:"lat"`
		Lng          float64 `json:"lng"`
		Phone        string  `json:"phone"`
		Email        string  `json:"email"`
		Description  string  `json:"description"`
		Rating       float64 `json:"rating"`
		TotalReviews int     `json:"total_reviews"`
		ImageURL     string  `json:"image_url"`
		IsActive     bool    `json:"is_active"`
		OpeningTime  string  `json:"opening_time"`
		ClosingTime  string  `json:"closing_time"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    string  `json:"updated_at"`
		IsFavorited  bool    `json:"is_favorited"`
	}

	var salons []FavoriteSalon
	for rows.Next() {
		var s FavoriteSalon
		rows.Scan(&s.ID, &s.OwnerID, &s.Name, &s.Address, &s.City, &s.State, &s.ZipCode,
			&s.Lat, &s.Lng, &s.Phone, &s.Email, &s.Description, &s.Rating,
			&s.TotalReviews, &s.ImageURL, &s.IsActive, &s.OpeningTime, &s.ClosingTime,
			&s.CreatedAt, &s.UpdatedAt)
		s.IsFavorited = true
		salons = append(salons, s)
	}
	if salons == nil {
		salons = []FavoriteSalon{}
	}

	c.JSON(http.StatusOK, salons)
}
