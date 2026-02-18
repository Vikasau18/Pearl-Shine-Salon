package handlers

import (
	"context"
	"net/http"
	"time"

	"saloon-backend/middleware"
	"saloon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default role
	if req.Role == "" {
		req.Role = "customer"
	}
	if req.Role != "customer" && req.Role != "salon_owner" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'customer' or 'salon_owner'"})
		return
	}

	// Check if email already exists
	var exists bool
	h.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Insert user
	var user models.User
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO users (email, password_hash, name, phone, role) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, email, name, phone, role, loyalty_points, created_at, updated_at`,
		req.Email, string(hashedPassword), req.Name, req.Phone, req.Role,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Phone, &user.Role, &user.LoyaltyPoints, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT
	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, email, password_hash, name, phone, role, COALESCE(avatar_url, ''), loyalty_points, created_at, updated_at 
		 FROM users WHERE email = $1`, req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.Role, &user.AvatarURL, &user.LoyaltyPoints, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.User
	err := h.DB.QueryRow(context.Background(),
		`SELECT id, email, name, phone, role, COALESCE(avatar_url, ''), loyalty_points, created_at, updated_at 
		 FROM users WHERE id = $1`, userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Phone, &user.Role, &user.AvatarURL, &user.LoyaltyPoints, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.DB.QueryRow(context.Background(),
		`UPDATE users SET name = COALESCE(NULLIF($1, ''), name), 
		 phone = COALESCE(NULLIF($2, ''), phone), 
		 avatar_url = COALESCE(NULLIF($3, ''), avatar_url),
		 updated_at = NOW()
		 WHERE id = $4
		 RETURNING id, email, name, phone, role, COALESCE(avatar_url, ''), loyalty_points, created_at, updated_at`,
		req.Name, req.Phone, req.AvatarURL, userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Phone, &user.Role, &user.AvatarURL, &user.LoyaltyPoints, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func generateToken(user models.User) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTSecret)
}
