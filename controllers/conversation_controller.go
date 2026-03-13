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

// CreateConversation handles POST /conversations - create or get a direct conversation
// @Summary Create direct conversation
// @Description Create a direct conversation between two users or return an existing one
// @Tags conversations
// @Accept json
// @Produce json
// @Param request body models.CreateConversationRequest true "Conversation members"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /conversations [post]
func (cc *ConversationController) CreateConversation(c *gin.Context) {
	var req models.CreateConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	conversation, err := cc.conversationService.CreateDirectConversation(req.Members)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Failed to create conversation: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(conversation, "Conversation ready"))
}
