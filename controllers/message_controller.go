package controllers

import (
	"chatapp-backend/models"
	"chatapp-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// MessageController handles message-related HTTP requests
type MessageController struct {
	messageService      *services.MessageService
	conversationService *services.ConversationService
	userService         *services.UserService
}

// NewMessageController creates a new message controller
func NewMessageController(msgService *services.MessageService, convService *services.ConversationService, userService *services.UserService) *MessageController {
	return &MessageController{
		messageService:      msgService,
		conversationService: convService,
		userService:         userService,
	}
}

// SendMessage handles POST /messages - send encrypted message
// @Summary Send encrypted message
// @Description Store an encrypted message without decryption
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.CreateMessageRequest true "Message data"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /messages [post]
func (mc *MessageController) SendMessage(c *gin.Context) {
	var req models.CreateMessageRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	// Extract user information from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated"))
		return
	}

	claimsMap := claims.(jwt.MapClaims)
	auth0ID := claimsMap["sub"].(string)

	// Extract email and username from claims
	email := ""
	if emailClaim, ok := claimsMap["email"].(string); ok {
		email = emailClaim
	}

	username := ""
	if nameClaim, ok := claimsMap["name"].(string); ok {
		username = nameClaim
	} else if nicknameClaim, ok := claimsMap["nickname"].(string); ok {
		username = nicknameClaim
	} else {
		username = auth0ID // Fallback to auth0 ID
	}

	// Create or update user in database
	user, err := mc.userService.GetOrCreateUser(auth0ID, username, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to process user: "+err.Error()))
		return
	}

	// Override senderID with authenticated user's ID
	req.SenderID = auth0ID

	// Ensure conversation exists and get the actual conversation ID
	actualConvID, err := mc.conversationService.EnsureConversationExists(req.ConversationID, req.SenderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create conversation: "+err.Error()))
		return
	}

	// Use the actual conversation ID for the message
	req.ConversationID = actualConvID

	// Create message
	message, err := mc.messageService.CreateMessage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create message: "+err.Error()))
		return
	}

	// Update conversation timestamp
	mc.conversationService.UpdateConversationTimestamp(actualConvID, message.Timestamp, message.EncryptedText)

	// Return response with user info and actual conversation ID used
	responseData := map[string]interface{}{
		"message":        message,
		"conversationId": actualConvID, // Return the actual conversation ID used
		"user": map[string]interface{}{
			"id":       user.ID,
			"auth0Id":  user.Auth0ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(responseData, "Message sent successfully"))
}

// GetMessages handles GET /messages/:conversationId - get all messages for a conversation
// @Summary Get messages by conversation
// @Description Retrieve all encrypted messages for a conversation
// @Tags messages
// @Produce json
// @Param conversationId path string true "Conversation ID"
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /messages/{conversationId} [get]
func (mc *MessageController) GetMessages(c *gin.Context) {
	conversationID := c.Param("conversationId")

	if conversationID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Conversation ID is required"))
		return
	}

	messages, err := mc.messageService.GetMessagesByConversation(conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to retrieve messages: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(messages, "Messages retrieved successfully"))
}
