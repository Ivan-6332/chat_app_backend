package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Conversation represents a conversation between users
type Conversation struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Members       []string           `json:"members" bson:"members"` // Changed from Participants to Members for consistency
	LastMessage   string             `json:"lastMessage,omitempty" bson:"last_message,omitempty"`
	LastMessageAt time.Time          `json:"lastMessageAt" bson:"last_message_at"`
	CreatedAt     time.Time          `json:"createdAt" bson:"created_at"`
}

// CreateConversationRequest represents the request body for creating a direct conversation
type CreateConversationRequest struct {
	// Preferred flow: send only contactUserId, backend will pair it with authenticated user.
	ContactUserID string `json:"contactUserId,omitempty"`
	// Backward compatibility: frontend may still send both members.
	Members []string `json:"members,omitempty"`
}
