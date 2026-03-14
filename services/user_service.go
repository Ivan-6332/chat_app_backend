package services

import (
	"chatapp-backend/config"
	"chatapp-backend/database"
	"chatapp-backend/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService handles user business logic with MongoDB
type UserService struct {
	collection *mongo.Collection
}

type auth0MgmtTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type auth0User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// Auth0SyncSummary provides counters for a sync run.
type Auth0SyncSummary struct {
	Fetched  int `json:"fetched"`
	Upserted int `json:"upserted"`
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		collection: database.GetCollection("users"),
	}
}

// GetOrCreateUser gets an existing user or creates a new one from Auth0 claims
func (s *UserService) GetOrCreateUser(auth0ID, username, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to find existing user
	filter := bson.M{"auth0_id": auth0ID}

	var user models.User
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err == nil {
		updateSet := bson.M{
			"last_seen_at": time.Now(),
		}

		if username != "" && username != user.Username {
			updateSet["username"] = username
			user.Username = username
		}

		if email != "" && email != user.Email {
			updateSet["email"] = email
			user.Email = email
		}

		_, updateErr := s.collection.UpdateOne(ctx, filter, bson.M{"$set": updateSet})
		if updateErr != nil {
			return nil, updateErr
		}

		user.LastSeenAt = updateSet["last_seen_at"].(time.Time)
		return &user, nil
	}

	// If not found, create new user
	if err == mongo.ErrNoDocuments {
		user = models.User{
			Auth0ID:    auth0ID,
			Username:   username,
			Email:      email,
			CreatedAt:  time.Now(),
			LastSeenAt: time.Now(),
		}

		_, err := s.collection.InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, err
}

// UpdateLastSeen updates the user's last seen timestamp
func (s *UserService) UpdateLastSeen(auth0ID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"auth0_id": auth0ID}
	update := bson.M{"$set": bson.M{"last_seen_at": time.Now()}}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetUserByAuth0ID gets a user by their Auth0 ID
func (s *UserService) GetUserByAuth0ID(auth0ID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"auth0_id": auth0ID}

	var user models.User
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SearchUsersByUsername searches users by username using a case-insensitive partial match
func (s *UserService) SearchUsersByUsername(usernameQuery string, excludeAuth0ID string, limit int64) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	normalized := strings.TrimSpace(usernameQuery)
	if normalized == "" {
		return []models.User{}, nil
	}

	quoted := regexp.QuoteMeta(normalized)

	filter := bson.M{
		"$or": []bson.M{
			{
				"username": bson.M{
					"$regex":   quoted,
					"$options": "i",
				},
			},
			{
				"display_name": bson.M{
					"$regex":   quoted,
					"$options": "i",
				},
			},
			{
				"email": bson.M{
					"$regex":   quoted,
					"$options": "i",
				},
			},
		},
	}

	if excludeAuth0ID != "" {
		filter["auth0_id"] = bson.M{"$ne": excludeAuth0ID}
	}

	opts := options.Find().SetLimit(limit)

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	if users == nil {
		return []models.User{}, nil
	}

	return users, nil
}

// SyncUsersFromAuth0 imports users from Auth0 Management API and upserts into MongoDB.
func (s *UserService) SyncUsersFromAuth0(maxPages int) (*Auth0SyncSummary, error) {
	if config.AppConfig.Auth0ClientID == "" || config.AppConfig.Auth0ClientSecret == "" {
		return nil, fmt.Errorf("AUTH0_CLIENT_ID and AUTH0_CLIENT_SECRET are required for sync")
	}

	if maxPages <= 0 {
		maxPages = 1
	}

	accessToken, err := s.getAuth0ManagementAccessToken()
	if err != nil {
		return nil, err
	}

	summary := &Auth0SyncSummary{}
	pageSize := config.AppConfig.Auth0SyncPageSize
	if pageSize <= 0 {
		pageSize = 100
	}

	for page := 0; page < maxPages; page++ {
		users, err := s.fetchAuth0Users(accessToken, page, pageSize)
		if err != nil {
			return nil, err
		}

		if len(users) == 0 {
			break
		}

		summary.Fetched += len(users)
		upserted, err := s.upsertAuth0Users(users)
		if err != nil {
			return nil, err
		}
		summary.Upserted += upserted

		if len(users) < pageSize {
			break
		}
	}

	return summary, nil
}

func (s *UserService) getAuth0ManagementAccessToken() (string, error) {
	endpoint := fmt.Sprintf("https://%s/oauth/token", config.AppConfig.Auth0Domain)

	payload := map[string]string{
		"client_id":     config.AppConfig.Auth0ClientID,
		"client_secret": config.AppConfig.Auth0ClientSecret,
		"audience":      config.AppConfig.Auth0ManagementAudience,
		"grant_type":    "client_credentials",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get Auth0 management token: status %d body %s", resp.StatusCode, string(respBytes))
	}

	var tokenResp auth0MgmtTokenResponse
	if err := json.Unmarshal(respBytes, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("Auth0 management token response missing access_token")
	}

	return tokenResp.AccessToken, nil
}

func (s *UserService) fetchAuth0Users(accessToken string, page int, pageSize int) ([]auth0User, error) {
	baseURL := fmt.Sprintf("https://%s/api/v2/users", config.AppConfig.Auth0Domain)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("per_page", fmt.Sprintf("%d", pageSize))
	q.Set("include_totals", "false")
	q.Set("fields", "user_id,username,nickname,name,email")
	q.Set("search_engine", "v3")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch Auth0 users: status %d body %s", resp.StatusCode, string(respBytes))
	}

	var users []auth0User
	if err := json.Unmarshal(respBytes, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) upsertAuth0Users(users []auth0User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	count := 0

	for _, u := range users {
		if u.UserID == "" {
			continue
		}

		username := u.Username
		if username == "" {
			if u.Nickname != "" {
				username = u.Nickname
			} else if u.Name != "" {
				username = u.Name
			} else {
				username = u.UserID
			}
		}

		filter := bson.M{"auth0_id": u.UserID}
		update := bson.M{
			"$set": bson.M{
				"auth0_id":     u.UserID,
				"username":     username,
				"email":        u.Email,
				"display_name": u.Name,
				"last_seen_at": now,
			},
			"$setOnInsert": bson.M{
				"created_at": now,
			},
		}

		_, err := s.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}
