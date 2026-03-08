# Quick Start Guide

## 🚀 Getting Started in 5 Minutes

### Step 1: Configure Environment

Copy the example environment file and add your Auth0 credentials:

```bash
# Copy the example file
cp .env.example .env

# Edit .env and add your Auth0 credentials
# AUTH0_DOMAIN=your-tenant.auth0.com
# AUTH0_AUDIENCE=your-api-identifier
```

### Step 2: Get Auth0 Credentials

1. Go to [Auth0 Dashboard](https://manage.auth0.com/)
2. Navigate to **Applications** → **APIs**
3. Create a new API or select existing one
4. Copy the **Identifier** (this is your `AUTH0_AUDIENCE`)
5. Your **Domain** is shown in the API settings (e.g., `dev-xxx.auth0.com`)

### Step 3: Install & Run

**Option A: Using Setup Script (Recommended)**

Windows (PowerShell):
```powershell
.\setup.ps1
.\chatapp-backend.exe
```

Linux/Mac:
```bash
chmod +x setup.sh
./setup.sh
./chatapp-backend
```

**Option B: Manual Setup**

```bash
# Install dependencies
go mod download

# Run the server
go run main.go
```

### Step 4: Test the API

The server will start on `http://localhost:8080`

**Test health endpoint (no auth required):**
```bash
curl http://localhost:8080/health
```

You should see:
```json
{
  "success": true,
  "data": {
    "service": "chatapp-backend",
    "status": "healthy"
  },
  "message": "Service is running"
}
```

### Step 5: Get an Auth Token

To test protected endpoints, you need a JWT token from Auth0:

1. **Frontend Integration**: Your Flutter app should handle authentication
2. **Testing**: Use Auth0's test token or implement a simple auth flow

Example using Auth0's test endpoint:
```bash
curl --request POST \
  --url https://YOUR_DOMAIN/oauth/token \
  --header 'content-type: application/json' \
  --data '{
    "client_id":"YOUR_CLIENT_ID",
    "client_secret":"YOUR_CLIENT_SECRET",
    "audience":"YOUR_API_IDENTIFIER",
    "grant_type":"client_credentials"
  }'
```

### Step 6: Make API Calls

**Send a message:**
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "senderId": "user123",
    "conversationId": "conv456",
    "encryptedText": "your_encrypted_message"
  }'
```

**Get messages:**
```bash
curl http://localhost:8080/api/v1/messages/conv456 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get conversations:**
```bash
curl http://localhost:8080/api/v1/conversations/user123 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 📱 Connecting Your Flutter App

Update your Flutter app to point to:
```dart
final baseUrl = 'http://localhost:8080/api/v1';
```

For Android emulator, use: `http://10.0.2.2:8080/api/v1`

For iOS simulator, use: `http://localhost:8080/api/v1`

For physical device, use your computer's IP: `http://192.168.x.x:8080/api/v1`

## 🐛 Troubleshooting

**"Authorization header is required"**
- Make sure you're including the Bearer token in the Authorization header
- Check token format: `Authorization: Bearer <token>`

**"Invalid token: invalid audience"**
- Verify `AUTH0_AUDIENCE` in .env matches your Auth0 API identifier

**"unable to find appropriate key"**
- Check `AUTH0_DOMAIN` in .env is correct
- Ensure the domain doesn't include `https://` or trailing `/`

**CORS errors from Flutter app**
- Current config allows all origins (development only)
- For production, update CORS settings in `main.go`

## 📚 Next Steps

- [ ] Integrate with a database (PostgreSQL, MongoDB)
- [ ] Add WebSocket support for real-time messaging
- [ ] Implement message pagination
- [ ] Add file/image upload endpoints
- [ ] Set up proper logging and monitoring
- [ ] Deploy to cloud (AWS, GCP, Azure)

## 🎯 Postman Collection

Import `postman_collection.json` into Postman for easy API testing!
