package main

import (
	"JWT/internal/handlers"
	"JWT/internal/middlewares"
	"JWT/internal/repository"
	"JWT/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	// Initialize Gin
	r := gin.Default()

	// Ping route (health check)
	r.GET("/ping", PingHandler)

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Dependency injection
	repo := repository.NewInMemory()
	jwtService := &services.JWTService{}
	handler := handlers.NewHandler(repo, jwtService)

	// API v1 routes
	v1 := r.Group("/api/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.POST("/signup", handler.Signup)
	auth.POST("/login", handler.Login)

	// User routes (secured)
	user := v1.Group("/user")
	user.GET("/getUsers", middlewares.AuthorizationMiddleware(), handler.GetAllUsers)

	// Start server
	if err := r.Run("localhost:8090"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
