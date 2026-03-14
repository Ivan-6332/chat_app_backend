package routes

import (
	"chatapp-backend/controllers"
	"chatapp-backend/middleware"
	"chatapp-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	messageController *controllers.MessageController,
	conversationController *controllers.ConversationController,
	userController *controllers.UserController,
) {
	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.SuccessResponse(gin.H{
			"status":  "healthy",
			"service": "chatapp-backend",
		}, "Service is running"))
	})

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Apply Auth0 middleware to all API routes
	v1.Use(middleware.Auth0Middleware())
	{
		// Message routes
		v1.POST("/messages", messageController.SendMessage)
		v1.GET("/messages/:conversationId", messageController.GetMessages)

		// Conversation routes
		v1.POST("/conversations", conversationController.CreateConversation)
		v1.GET("/conversations/:userId", conversationController.GetConversations)

		// User routes
		v1.GET("/users/me", userController.GetCurrentUser)
		v1.GET("/users/search", userController.SearchUsers)
		v1.POST("/users/sync", userController.SyncUsers)
		v1.GET("/users/:id", userController.GetUserByID)
	}
}
