# ChatApp Backend

Go backend REST API for a secure encrypted chat application.

## Features

- 🔐 JWT authentication via Auth0
- 📨 Encrypted message storage and retrieval
- 💬 Conversation management
- 🚀 Built with Gin framework
- ✨ Clean architecture (routes, controllers, services)

## Architecture

```
chatapp-backend/
├── main.go                 # Application entry point
├── config/                 # Configuration management
├── middleware/             # Auth0 JWT middleware
├── models/                 # Data models
├── controllers/            # HTTP request handlers
├── services/               # Business logic
└── routes/                 # Route definitions
```

## API Endpoints

All endpoints (except `/health`) require a valid Auth0 JWT token in the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

### Health Check
- `GET /health` - Check service status (no auth required)

### Messages
- `POST /api/v1/messages` - Send encrypted message
  ```json
  {
    "senderId": "user-id",
    "conversationId": "conversation-id",
    "encryptedText": "encrypted-message-content"
  }
  ```

- `GET /api/v1/messages/:conversationId` - Get all messages for a conversation

### Conversations
- `GET /api/v1/conversations/:userId` - Get all conversations for a user

## Setup

### Prerequisites
- Go 1.21 or higher
- Auth0 account with configured API

### Environment Variables

Create a `.env` file in the root directory:

```env
AUTH0_DOMAIN=your-tenant.auth0.com
AUTH0_AUDIENCE=your-api-identifier
PORT=8080
GIN_MODE=debug
```

### Installation

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Development

### Project Structure

- **config**: Application configuration and environment variable management
- **middleware**: Auth0 JWT token validation
- **models**: Data structures for messages, conversations, and API responses
- **services**: Business logic layer with in-memory storage
- **controllers**: HTTP request handlers
- **routes**: API route definitions

### Security

- All API endpoints are protected with Auth0 JWT middleware
- Messages are stored in encrypted form - backend never decrypts them
- JWKS caching for efficient token validation
- CORS configuration for cross-origin requests

## Production Considerations

⚠️ **Current Implementation Notes:**

- Uses in-memory storage (data is lost on restart)
- For production, implement:
  - Database integration (PostgreSQL, MongoDB, etc.)
  - Redis for caching
  - Rate limiting
  - Request logging
  - Error tracking
  - Health monitoring
  - Proper CORS configuration (restrict origins)

## Testing

Test the API with curl:

```bash
# Health check
curl http://localhost:8080/health

# Send message (with auth token)
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "senderId": "user123",
    "conversationId": "conv456",
    "encryptedText": "encrypted_message_here"
  }'

# Get messages
curl http://localhost:8080/api/v1/messages/conv456 \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get conversations
curl http://localhost:8080/api/v1/conversations/user123 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## License

MIT
