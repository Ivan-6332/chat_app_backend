package main

import (
	"chatapp-backend/config"
	"chatapp-backend/controllers"
	"chatapp-backend/database"
	"chatapp-backend/routes"
	"chatapp-backend/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to MongoDB
	if err := database.ConnectMongoDB(config.AppConfig.MongoDBURI, config.AppConfig.MongoDBName); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.DisconnectMongoDB()

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Initialize services
	messageService := services.NewMessageService()
	conversationService := services.NewConversationService(messageService)
	userService := services.NewUserService()

	// Initialize controllers
	messageController := controllers.NewMessageController(messageService, conversationService, userService)
	conversationController := controllers.NewConversationController(conversationService)
	userController := controllers.NewUserController(userService)

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // In production, specify your Flutter app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(router, messageController, conversationController, userController)

	// Create HTTP server
	port := ":" + config.AppConfig.Port
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		log.Printf("Auth0 Domain: %s", config.AppConfig.Auth0Domain)
		log.Printf("Auth0 Audience: %s", config.AppConfig.Auth0Audience)
		log.Printf("MongoDB Database: %s", config.AppConfig.MongoDBName)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
