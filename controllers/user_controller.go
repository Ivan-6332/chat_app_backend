package controllers

import (
	"chatapp-backend/config"
	"chatapp-backend/services"
	"net/http"
	"strconv"
	"strings"

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

	email := ""
	if emailClaim, ok := claimsMap["email"].(string); ok {
		email = emailClaim
	}

	username := ""
	if preferredUsername, ok := claimsMap["preferred_username"].(string); ok && preferredUsername != "" {
		username = preferredUsername
	} else if nicknameClaim, ok := claimsMap["nickname"].(string); ok && nicknameClaim != "" {
		username = nicknameClaim
	} else if nameClaim, ok := claimsMap["name"].(string); ok && nameClaim != "" {
		username = nameClaim
	} else {
		username = auth0ID
	}

	// Ensure user exists in local DB so contacts/search can operate on Mongo-backed user data.
	user, err := uc.userService.GetOrCreateUser(auth0ID, username, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to process user profile"))
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

	// Defensive fallback: if route matching sends /users/search here,
	// delegate to the search handler instead of returning a misleading 404.
	if strings.HasPrefix(strings.ToLower(auth0ID), "search") {
		uc.SearchUsers(c)
		return
	}

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

	includeSelf := false
	if c.Query("includeSelf") == "true" {
		includeSelf = true
	}

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated"))
		return
	}

	claimsMap := claims.(jwt.MapClaims)
	auth0ID := claimsMap["sub"].(string)
	excludeAuth0ID := auth0ID
	if includeSelf {
		excludeAuth0ID = ""
	}

	users, err := uc.userService.SearchUsersByUsername(usernameQuery, excludeAuth0ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to search users: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(users, "Users retrieved successfully"))
}

// SyncUsers handles POST /users/sync - synchronize Auth0 users into MongoDB
// @Summary Sync users from Auth0
// @Description Pull users from Auth0 Management API and upsert into MongoDB
// @Tags users
// @Produce json
// @Param maxPages query int false "Max pages to sync (default 3)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 403 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /users/sync [post]
func (uc *UserController) SyncUsers(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("User not authenticated"))
		return
	}

	claimsMap := claims.(jwt.MapClaims)
	if !hasScope(claimsMap, "sync:users") {
		c.JSON(http.StatusForbidden, models.ErrorResponse("Missing required scope: sync:users"))
		return
	}

	if config.AppConfig.Auth0SyncServiceTokenOnly && !isServiceToken(claimsMap) {
		c.JSON(http.StatusForbidden, models.ErrorResponse("/users/sync is restricted to service/admin tokens"))
		return
	}

	callerClientID := getCallerClientID(claimsMap)
	if len(config.AppConfig.Auth0SyncAllowedClientIDs) > 0 {
		if callerClientID == "" || !contains(config.AppConfig.Auth0SyncAllowedClientIDs, callerClientID) {
			c.JSON(http.StatusForbidden, models.ErrorResponse("caller client_id is not allowed for sync"))
			return
		}
	}

	maxPages := 3
	if raw := c.Query("maxPages"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse("maxPages must be a positive integer"))
			return
		}
		if parsed > 20 {
			parsed = 20
		}
		maxPages = parsed
	}

	summary, err := uc.userService.SyncUsersFromAuth0(maxPages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to sync users: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(summary, "Users synchronized successfully"))
}

func hasScope(claims jwt.MapClaims, required string) bool {
	if permissionsRaw, ok := claims["permissions"]; ok {
		if permissions, ok := permissionsRaw.([]interface{}); ok {
			for _, p := range permissions {
				if scope, ok := p.(string); ok && scope == required {
					return true
				}
			}
		}
	}

	if scopeRaw, ok := claims["scope"]; ok {
		if scopeText, ok := scopeRaw.(string); ok {
			for _, scope := range strings.Fields(scopeText) {
				if scope == required {
					return true
				}
			}
		}
	}

	return false
}

func isServiceToken(claims jwt.MapClaims) bool {
	if gtyRaw, ok := claims["gty"]; ok {
		if gty, ok := gtyRaw.(string); ok && strings.EqualFold(gty, "client-credentials") {
			return true
		}
	}

	if subRaw, ok := claims["sub"]; ok {
		if sub, ok := subRaw.(string); ok {
			return strings.HasSuffix(sub, "@clients")
		}
	}

	return false
}

func getCallerClientID(claims jwt.MapClaims) string {
	if azpRaw, ok := claims["azp"]; ok {
		if azp, ok := azpRaw.(string); ok {
			return azp
		}
	}

	if clientIDRaw, ok := claims["client_id"]; ok {
		if clientID, ok := clientIDRaw.(string); ok {
			return clientID
		}
	}

	return ""
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}

	return false
}
