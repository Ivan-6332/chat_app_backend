package services

import (
	"chatapp-backend/database"
	"chatapp-backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConversationService handles conversation business logic with MongoDB
type ConversationService struct {
	collection *mongo.Collection
	msgService *MessageService
}

// NewConversationService creates a new conversation service
func NewConversationService(msgService *MessageService) *ConversationService {
	return &ConversationService{
		collection: database.GetCollection("conversations"),
		msgService: msgService,
	}
}

// GetOrCreateConversation gets an existing conversation or creates a new one
func (s *ConversationService) GetOrCreateConversation(conversationID string, members []string) (*models.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to find existing conversation with these members
	filter := bson.M{"members": bson.M{"$all": members}}

	var conversation models.Conversation
	err := s.collection.FindOne(ctx, filter).Decode(&conversation)
	if err == nil {
		return &conversation, nil
	}

	// If not found, create new conversation
	if err == mongo.ErrNoDocuments {
		conversation = models.Conversation{
			ID:            primitive.NewObjectID(),
			Members:       members,
			CreatedAt:     time.Now(),
			LastMessageAt: time.Now(),
		}

		_, err := s.collection.InsertOne(ctx, conversation)
		if err != nil {
			return nil, err
		}

		return &conversation, nil
	}

	return nil, err
}

// GetConversationsByUser retrieves all conversations for a user from MongoDB
func (s *ConversationService) GetConversationsByUser(userID string) ([]models.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find conversations where user is a member
	filter := bson.M{"members": userID}
	opts := options.Find().SetSort(bson.D{{Key: "last_message_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}

	// Populate last message for each conversation
	for i := range conversations {
		convIDStr := conversations[i].ID.Hex()
		if msg, err := s.msgService.GetLatestMessage(convIDStr); err == nil {
			conversations[i].LastMessage = msg.EncryptedText
			conversations[i].LastMessageAt = msg.Timestamp
		}
	}

	if conversations == nil {
		return []models.Conversation{}, nil
	}

	return conversations, nil
}

// UpdateConversationTimestamp updates the last message timestamp of a conversation
func (s *ConversationService) UpdateConversationTimestamp(conversationID string, timestamp time.Time, lastMessage string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"last_message_at": timestamp,
			"last_message":    lastMessage,
		},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	return err
}

// EnsureConversationExists creates a conversation if it doesn't exist and returns the conversation ID
func (s *ConversationService) EnsureConversationExists(conversationID string, senderID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		// If invalid ObjectID, create a new conversation
		newID := primitive.NewObjectID()
		conversation := models.Conversation{
			ID:            newID,
			Members:       []string{senderID},
			CreatedAt:     time.Now(),
			LastMessageAt: time.Now(),
		}

		_, err := s.collection.InsertOne(ctx, conversation)
		if err != nil {
			return "", err
		}
		return newID.Hex(), nil
	}

	// Check if conversation exists
	filter := bson.M{"_id": objID}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return "", err
	}

	if count > 0 {
		return objID.Hex(), nil // Conversation exists
	}

	// Create new conversation
	conversation := models.Conversation{
		ID:            objID,
		Members:       []string{senderID},
		CreatedAt:     time.Now(),
		LastMessageAt: time.Now(),
	}

	_, err = s.collection.InsertOne(ctx, conversation)
	if err != nil {
		return "", err
	}
	return objID.Hex(), nil
}
