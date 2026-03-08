package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(50).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	Database = client.Database(dbName)

	log.Println("✅ Connected to MongoDB successfully")

	// Create indexes
	if err := createIndexes(); err != nil {
		log.Printf("⚠️ Warning: Failed to create indexes: %v", err)
	}

	return nil
}

// createIndexes creates necessary indexes for performance
func createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Index for messages collection
	messagesCollection := Database.Collection("messages")
	_, err := messagesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "conversation_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "conversation_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
	})
	if err != nil {
		return err
	}

	// Index for conversations collection
	conversationsCollection := Database.Collection("conversations")
	_, err = conversationsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "members", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "last_message_at", Value: -1}},
		},
	})
	if err != nil {
		return err
	}

	// Index for users collection
	usersCollection := Database.Collection("users")
	_, err = usersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "auth0_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	log.Println("✅ Database indexes created successfully")
	return nil
}

// DisconnectMongoDB closes the MongoDB connection
func DisconnectMongoDB() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return Client.Disconnect(ctx)
	}
	return nil
}

// GetCollection returns a collection from the database
func GetCollection(name string) *mongo.Collection {
	return Database.Collection(name)
}
