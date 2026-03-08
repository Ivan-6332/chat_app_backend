package services

import (
	"chatapp-backend/database"
	"chatapp-backend/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MessageService handles message business logic with MongoDB
type MessageService struct {
	collection *mongo.Collection
}

// NewMessageService creates a new message service
func NewMessageService() *MessageService {
	return &MessageService{
		collection: database.GetCollection("messages"),
	}
}

// CreateMessage creates a new encrypted message in MongoDB
func (s *MessageService) CreateMessage(req models.CreateMessageRequest) (*models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := models.Message{
		ID:             primitive.NewObjectID(),
		SenderID:       req.SenderID,
		ConversationID: req.ConversationID,
		EncryptedText:  req.EncryptedText, // Stored encrypted
		Timestamp:      time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// GetMessagesByConversation retrieves all messages for a conversation from MongoDB
func (s *MessageService) GetMessagesByConversation(conversationID string) ([]models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find messages for this conversation, sorted by timestamp
	filter := bson.M{"conversation_id": conversationID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	if messages == nil {
		return []models.Message{}, nil // Return empty array if no messages
	}

	return messages, nil
}

// GetLatestMessage gets the most recent message for a conversation
func (s *MessageService) GetLatestMessage(conversationID string) (*models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"conversation_id": conversationID}
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	var message models.Message
	err := s.collection.FindOne(ctx, filter, opts).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no messages found")
		}
		return nil, err
	}

	return &message, nil
}
