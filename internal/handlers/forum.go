package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naufalwaiz/hukum-chatbot/internal/services"
)

// Request body
type ForumRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

// Create
func CreateForum(c *gin.Context) {
	var req ForumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	forum, err := services.CreateForum(req.Title, req.Body, req.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "news created successfully",
		"id":      forum.ID,
		"author":  forum.Author,
	})
}

// Get all
func GetAllForums(c *gin.Context) {
	forums, err := services.GetAllForums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, forums)
}

// Get by ID
func GetForumByID(c *gin.Context) {
	id := c.Param("id")
	forum, err := services.GetForumByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "forum not found"})
		return
	}
	c.JSON(http.StatusOK, forum)
}

// Update
func UpdateForum(c *gin.Context) {
	id := c.Param("id")
	var req ForumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateForum(id, req.Title, req.Body, req.Author); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "forum updated"})
}

// Delete
func DeleteForum(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteForum(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "forum deleted"})
}
