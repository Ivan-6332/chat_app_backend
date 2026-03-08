package services

import (
	"chatapp-backend/database"
	"chatapp-backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService handles user business logic with MongoDB
type UserService struct {
	collection *mongo.Collection
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
		// Update last seen
		s.UpdateLastSeen(auth0ID)
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
