package handlers

import (
    "log"
    "database/sql"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/DomT00T/actico-auth-api/internal/auth"
    "github.com/DomT00T/actico-auth-api/internal/models"
    "github.com/DomT00T/actico-auth-api/pkg/utils"
)


// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := utils.HashPassword(user.PasswordHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		now := time.Now()
		_, err = db.Exec(`
			INSERT INTO users (email, password_hash, user_type, first_name, last_name, phone_number, gender, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, user.Email, hashedPassword, user.UserType, user.FirstName, user.LastName, user.PhoneNumber, user.Gender, now, now)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", loginData.Email).Scan(&user.ID, &user.PasswordHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if !utils.CheckPasswordHash(loginData.Password, user.PasswordHash) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := auth.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		_, err = db.Exec("UPDATE users SET last_login = $1 WHERE id = $2", time.Now(), user.ID)
		if err != nil {
			// Log the error, but don't return it to the user
			log.Printf("Error updating last login: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// Implement GoogleLogin, GoogleCallback, FacebookLogin, and FacebookCallback handlers here