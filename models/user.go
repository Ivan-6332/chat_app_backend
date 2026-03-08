package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Auth0ID     string             `json:"auth0Id" bson:"auth0_id"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
	DisplayName string             `json:"displayName,omitempty" bson:"display_name,omitempty"`
	ProfilePic  string             `json:"profilePic,omitempty" bson:"profile_pic,omitempty"`
	PublicKey   string             `json:"publicKey,omitempty" bson:"public_key,omitempty"` // For E2E encryption
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
	LastSeenAt  time.Time          `json:"lastSeenAt" bson:"last_seen_at"`
}
