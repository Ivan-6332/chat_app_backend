package controllers

import (
	"chatapp-backend/services"
	"net/http"

	"chatapp-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetCurrentUser handles GET /users/me - get current authenticated user
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Produce json
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /users/me [get]
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	// Extract user information from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated"))
		return
	}

	claimsMap := claims.(jwt.MapClaims)
	auth0ID := claimsMap["sub"].(string)

	// Get user from database
	user, err := uc.userService.GetUserByAuth0ID(auth0ID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(user, "User retrieved successfully"))
}

// GetUserByID handles GET /users/:id - get user by ID
// @Summary Get user by ID
// @Description Get a user's profile by their Auth0 ID
// @Tags users
// @Produce json
// @Param id path string true "User Auth0 ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(c *gin.Context) {
	auth0ID := c.Param("id")

	if auth0ID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("User ID is required"))
		return
	}

	user, err := uc.userService.GetUserByAuth0ID(auth0ID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(user, "User retrieved successfully"))
}
