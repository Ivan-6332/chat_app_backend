package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents an encrypted message in the system
type Message struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SenderID       string             `json:"senderId" bson:"sender_id"`
	ConversationID string             `json:"conversationId" bson:"conversation_id"`
	EncryptedText  string             `json:"encryptedText" bson:"encrypted_text"` // Encrypted in DB
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
}

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	SenderID       string `json:"senderId" binding:"required"`
	ConversationID string `json:"conversationId" binding:"required"`
	EncryptedText  string `json:"encryptedText" binding:"required"`
}
