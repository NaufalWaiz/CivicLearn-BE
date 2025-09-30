package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/naufalwaiz/hukum-chatbot/config"
	"github.com/naufalwaiz/hukum-chatbot/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	id := uuid.New()

	_, err := config.DB.Exec(
		"INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
		id, req.Email, string(hashed),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id, password string
	err := config.DB.QueryRow(
		"SELECT id, password FROM users WHERE email=$1", req.Email,
	).Scan(&id, &password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := services.GenerateToken(id)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
