package controllers

import (
	"chatapp-backend/models"
	"chatapp-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ConversationController handles conversation-related HTTP requests
type ConversationController struct {
	conversationService *services.ConversationService
}

// NewConversationController creates a new conversation controller
func NewConversationController(convService *services.ConversationService) *ConversationController {
	return &ConversationController{
		conversationService: convService,
	}
}

// GetConversations handles GET /conversations/:userId - get all conversations for a user
// @Summary Get user conversations
// @Description Retrieve all conversations for a specific user
// @Tags conversations
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /conversations/{userId} [get]
func (cc *ConversationController) GetConversations(c *gin.Context) {
	userID := c.Param("userId")

	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("User ID is required"))
		return
	}

	conversations, err := cc.conversationService.GetConversationsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to retrieve conversations: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(conversations, "Conversations retrieved successfully"))
}
