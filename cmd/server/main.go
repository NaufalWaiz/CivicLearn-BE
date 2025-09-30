package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/naufalwaiz/hukum-chatbot/config"
	"github.com/naufalwaiz/hukum-chatbot/internal/handlers"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system env")
	}

	// Init DB
	config.InitDB()

	// Login and Register routes
	r := gin.Default()
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/chatbot", handlers.Chatbot)
	r.GET("/health", handlers.HealthCheck)

	// News routes
	r.POST("/news", handlers.CreateNews)
	r.GET("/news", handlers.GetAllNews)
	r.GET("/news/:id", handlers.GetNewsByID)
	r.PUT("/news/:id", handlers.UpdateNews)
	r.DELETE("/news/:id", handlers.DeleteNews)

	// Forum routes
	r.POST("/forum", handlers.CreateForum)
	r.GET("/forum", handlers.GetAllForums)
	r.GET("/forum/:id", handlers.GetForumByID)
	r.PUT("/forum/:id", handlers.UpdateForum)
	r.DELETE("/forum/:id", handlers.DeleteForum)

	r.Run(":8080")
}
