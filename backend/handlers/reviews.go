package handlers

import (
	"context"
	"net/http"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewHandler struct {
	DB *pgxpool.Pool
}

func NewReviewHandler(db *pgxpool.Pool) *ReviewHandler {
	return &ReviewHandler{DB: db}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := middleware.GetUserID(c)

	var staffID, appointmentID *string
	if req.StaffID != "" {
		staffID = &req.StaffID
	}
	if req.AppointmentID != "" {
		appointmentID = &req.AppointmentID
	}

	var review models.Review
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO reviews (customer_id, salon_id, staff_id, appointment_id, rating, comment)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, customer_id, salon_id, staff_id, appointment_id, rating, COALESCE(comment,''), created_at, updated_at`,
		customerID, req.SalonID, staffID, appointmentID, req.Rating, req.Comment,
	).Scan(&review.ID, &review.CustomerID, &review.SalonID, &review.StaffID,
		&review.AppointmentID, &review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review", "details": err.Error()})
		return
	}

	// Update salon average rating
	h.DB.Exec(context.Background(),
		`UPDATE salons SET 
		 rating = (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE salon_id = $1),
		 total_reviews = (SELECT COUNT(*) FROM reviews WHERE salon_id = $1),
		 updated_at = NOW()
		 WHERE id = $1`, req.SalonID)

	c.JSON(http.StatusCreated, review)
}

func (h *ReviewHandler) GetSalonReviews(c *gin.Context) {
	salonID := c.Param("id")

	rows, err := h.DB.Query(context.Background(),
		`SELECT r.id, r.customer_id, r.salon_id, r.staff_id, r.appointment_id, 
		 r.rating, COALESCE(r.comment,''), r.created_at, r.updated_at,
		 u.name, COALESCE(u.avatar_url,'')
		 FROM reviews r
		 JOIN users u ON u.id = r.customer_id
		 WHERE r.salon_id = $1
		 ORDER BY r.created_at DESC`, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		rows.Scan(&r.ID, &r.CustomerID, &r.SalonID, &r.StaffID, &r.AppointmentID,
			&r.Rating, &r.Comment, &r.CreatedAt, &r.UpdatedAt,
			&r.CustomerName, &r.CustomerAvatar)
		reviews = append(reviews, r)
	}
	if reviews == nil {
		reviews = []models.Review{}
	}

	c.JSON(http.StatusOK, reviews)
}
