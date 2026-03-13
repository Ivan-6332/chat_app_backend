package controllers

import (
	"chatapp-backend/services"
	"net/http"
	"strconv"

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

// SearchUsers handles GET /users/search?username={username}&limit={limit}
// @Summary Search users by username
// @Description Search users with case-insensitive partial match by username
// @Tags users
// @Produce json
// @Param username query string true "Username query"
// @Param limit query int false "Maximum results (default 20, max 50)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/search [get]
func (uc *UserController) SearchUsers(c *gin.Context) {
	usernameQuery := c.Query("username")
	if usernameQuery == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("username query parameter is required"))
		return
	}

	limit := int64(20)
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse("limit must be a positive integer"))
			return
		}

		if parsed > 50 {
			parsed = 50
		}

		limit = int64(parsed)
	}

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated"))
		return
	}

	claimsMap := claims.(jwt.MapClaims)
	auth0ID := claimsMap["sub"].(string)

	users, err := uc.userService.SearchUsersByUsername(usernameQuery, auth0ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to search users: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(users, "Users retrieved successfully"))
}
