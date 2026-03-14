package controllers

import (
	"chatapp-backend/models"
	"chatapp-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated"))
		return
	}

	claimsMap := claims.(jwt.MapClaims)
	auth0ID := claimsMap["sub"].(string)

	var members []string
	if req.ContactUserID != "" {
		if req.ContactUserID == auth0ID {
			c.JSON(http.StatusBadRequest, models.ErrorResponse("contactUserId must be a different user"))
			return
		}
		members = []string{auth0ID, req.ContactUserID}
	} else {
		if len(req.Members) != 2 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse("members must contain exactly 2 user IDs"))
			return
		}

		if req.Members[0] != auth0ID && req.Members[1] != auth0ID {
			c.JSON(http.StatusForbidden, models.ErrorResponse("authenticated user must be included in members"))
			return
		}

		members = req.Members
	}

	conversation, err := cc.conversationService.CreateDirectConversation(members)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Failed to create conversation: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(conversation, "Conversation ready"))
}
