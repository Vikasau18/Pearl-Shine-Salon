package handlers

import (
	"context"
	"net/http"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PromoHandler struct {
	DB *pgxpool.Pool
}

func NewPromoHandler(db *pgxpool.Pool) *PromoHandler {
	return &PromoHandler{DB: db}
}

func (h *PromoHandler) CreatePromo(c *gin.Context) {
	salonID := c.Param("salon_id")
	var req models.CreatePromoRequest
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

	var promo models.PromoCode
	err := h.DB.QueryRow(context.Background(),
		`INSERT INTO promo_codes (salon_id, code, discount_percent, valid_from, valid_until, max_uses)
		 VALUES ($1, $2, $3, COALESCE(NULLIF($4,'')::timestamptz, NOW()), NULLIF($5,'')::timestamptz, $6)
		 RETURNING id, salon_id, code, discount_percent, valid_from, valid_until, max_uses, used_count, is_active, created_at`,
		salonID, req.Code, req.DiscountPercent, req.ValidFrom, req.ValidUntil, req.MaxUses,
	).Scan(&promo.ID, &promo.SalonID, &promo.Code, &promo.DiscountPercent,
		&promo.ValidFrom, &promo.ValidUntil, &promo.MaxUses, &promo.UsedCount,
		&promo.IsActive, &promo.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create promo code", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, promo)
}

func (h *PromoHandler) GetSalonPromos(c *gin.Context) {
	salonID := c.Param("salon_id")

	rows, err := h.DB.Query(context.Background(),
		`SELECT id, salon_id, code, discount_percent, valid_from, valid_until, max_uses, used_count, is_active, created_at
		 FROM promo_codes WHERE salon_id = $1 ORDER BY created_at DESC`, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promo codes"})
		return
	}
	defer rows.Close()

	var promos []models.PromoCode
	for rows.Next() {
		var p models.PromoCode
		rows.Scan(&p.ID, &p.SalonID, &p.Code, &p.DiscountPercent, &p.ValidFrom,
			&p.ValidUntil, &p.MaxUses, &p.UsedCount, &p.IsActive, &p.CreatedAt)
		promos = append(promos, p)
	}
	if promos == nil {
		promos = []models.PromoCode{}
	}

	c.JSON(http.StatusOK, promos)
}

func (h *PromoHandler) ValidatePromo(c *gin.Context) {
	code := c.Query("code")
	salonID := c.Query("salon_id")

	if code == "" || salonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code and salon_id required"})
		return
	}

	var promo models.PromoCode
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, salon_id, code, discount_percent, valid_from, valid_until, max_uses, used_count, is_active
		 FROM promo_codes 
		 WHERE code = $1 AND salon_id = $2 AND is_active = true 
		 AND (valid_until IS NULL OR valid_until > NOW())`,
		code, salonID).Scan(&promo.ID, &promo.SalonID, &promo.Code, &promo.DiscountPercent,
		&promo.ValidFrom, &promo.ValidUntil, &promo.MaxUses, &promo.UsedCount, &promo.IsActive)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"valid": false, "error": "Invalid or expired promo code"})
		return
	}

	if promo.MaxUses != nil && promo.UsedCount >= *promo.MaxUses {
		c.JSON(http.StatusOK, gin.H{"valid": false, "error": "Promo code has reached maximum uses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":            true,
		"discount_percent": promo.DiscountPercent,
		"code":             promo.Code,
	})
}

func (h *PromoHandler) DeletePromo(c *gin.Context) {
	promoID := c.Param("promo_id")

	_, err := h.DB.Exec(context.Background(),
		"UPDATE promo_codes SET is_active = false WHERE id = $1", promoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate promo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promo code deactivated"})
}
