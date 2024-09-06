package main

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/DomT00T/actico-auth-api/cmd/api/docs"
    "github.com/DomT00T/actico-auth-api/internal/config"
    "github.com/DomT00T/actico-auth-api/internal/database"
    "github.com/DomT00T/actico-auth-api/internal/handlers"
    "github.com/DomT00T/actico-auth-api/internal/middlewares"
)


// @title Your API Title
// @version 1.0
// @description Your API Description
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create Gin router
	r := gin.Default()

	// Define API routes
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register(db))
			auth.POST("/login", handlers.Login(db))
			auth.GET("/google", handlers.GoogleLogin())
			auth.GET("/google/callback", handlers.GoogleCallback(db))
			auth.GET("/facebook", handlers.FacebookLogin())
			auth.GET("/facebook/callback", handlers.FacebookCallback(db))
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.GET("/user", handlers.GetUser(db))
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	log.Printf("Server is running on :%s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
