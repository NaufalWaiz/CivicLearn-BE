package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naufalwaiz/hukum-chatbot/internal/services"
)

// Request body
type NewsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

// Create
func CreateNews(c *gin.Context) {
	var req NewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, author, err := services.CreateNews(req.Title, req.Content, req.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "news created successfully",
		"id":      id,
		"author":  author,
	})
}

// Get all
func GetAllNews(c *gin.Context) {
	news, err := services.GetAllNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

// Get by ID
func GetNewsByID(c *gin.Context) {
	id := c.Param("id")
	news, err := services.GetNewsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		return
	}
	c.JSON(http.StatusOK, news)
}

// Update
func UpdateNews(c *gin.Context) {
	id := c.Param("id")
	var req NewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateNews(id, req.Title, req.Content, req.Author); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "news updated"})
}

// Delete
func DeleteNews(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteNews(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "news deleted"})
}
