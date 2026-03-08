# Quick Start Guide - MongoDB + Flutter Integration

## 🚀 Quick Setup (5 Minutes)

### 1️⃣ Create MongoDB Database
1. Go to https://www.mongodb.com/cloud/atlas/register
2. Sign up (FREE, no credit card)
3. Create cluster → Create user → Get connection string
4. **Detailed guide:** See [`MONGODB_SETUP.md`](MONGODB_SETUP.md)

### 2️⃣ Configure Backend
Edit `.env` file:
```env
MONGODB_URI=mongodb+srv://YOUR_USER:YOUR_PASS@cluster.mongodb.net/
```

### 3️⃣ Configure Flutter App
Edit `lib/services/api_service.dart`:
```dart
// Line 10-12
static const String baseUrl = 'http://10.0.2.2:8080'; // Android Emulator
// OR
static const String baseUrl = 'http://YOUR_IP:8080'; // Real Device
```

### 4️⃣ Start Backend
```bash
go run main.go
```

### 5️⃣ Run Flutter App
```bash
flutter run
```

---

## 📡 API Endpoints

All endpoints require `Authorization: Bearer <token>` header

### Messages
- **POST** `/api/v1/messages` - Send encrypted message
  ```json
  {
    "senderId": "auth0|123",
    "conversationId": "conv_456",
    "encryptedText": "U2FsdGVkX1..."
  }
  ```

- **GET** `/api/v1/messages/:conversationId` - Get messages

### Conversations
- **GET** `/api/v1/conversations/:userId` - Get user's conversations

### Health
- **GET** `/health` - Check server status (no auth required)

---

## 💻 Flutter Usage Examples

### Send a Message
```dart
import 'services/chat_service.dart';

await ChatService.sendMessage(
  plainText: 'Hello!',
  conversationId: 'conv_123',
  senderId: userId,
);
```

### Get Messages
```dart
final messages = await ChatService.getMessages(
  conversationId: 'conv_123',
);
// Messages are automatically decrypted!
```

### Get Conversations
```dart
final conversations = await ChatService.getConversations(
  userId: userId,
);
```

### Test Connection
```dart
final isConnected = await ChatService.testConnection();
print('Backend online: $isConnected');
```

---

## 🗂️ Project Structure

### Backend (Go)
```
chatapp-backend/
├── database/          # MongoDB connection
│   └── mongodb.go
├── models/           # Data structures
│   ├── user.go
│   ├── message.go
│   └── conversation.go
├── services/         # Business logic
│   ├── user_service.go
│   ├── message_service.go
│   └── conversation_service.go
├── controllers/      # HTTP handlers
├── middleware/       # Auth0 JWT validation
├── routes/          # Route definitions
└── main.go          # Entry point
```

### Flutter App
```
chatapp/lib/
├── services/
│   ├── auth_service.dart       # Auth0 login
│   ├── api_service.dart        # HTTP calls
│   ├── chat_service.dart       # Chat operations
│   └── encryption_service.dart # AES encryption
├── screens/         # UI screens
├── models/          # Data models
└── main.dart
```

---

## 🔍 Verify MongoDB Data

After sending messages:

1. Go to MongoDB Atlas Dashboard
2. Click **"Browse Collections"**
3. Select `chatapp_db` database
4. View collections:
   - `messages` - Encrypted messages
   - `conversations` - Chat conversations
   - `users` - User profiles

---

## ⚡ MongoDB Features Implemented

✅ **Connection Pooling** (10-50 connections)  
✅ **Automatic Indexes** for fast queries  
✅ **Encrypted Storage** (messages stored encrypted)  
✅ **Graceful Shutdown** handling  
✅ **Error Recovery** with timeouts  

---

## 🐛 Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| MongoDB connection failed | Check `MONGODB_URI` in `.env`, verify password |
| Port 8080 already in use | Kill process: `Get-Process -Name go \| Stop-Process` |
| Flutter can't connect | Use `10.0.2.2` for emulator, your IP for device |
| Decryption failed | Ensure same encryption key on all devices |

---

## 📊 Performance Optimizations

- ✅ Indexed `conversation_id` for fast message lookup
- ✅ Indexed `members` for fast conversation queries
- ✅ Sorted by `timestamp` for chronological order
- ✅ Connection pooling reduces overhead
- ✅ 5-10 second timeouts prevent hanging

---

## 🔐 Security Checklist

- ✅ Messages encrypted client-side (AES-256)
- ✅ Messages stored encrypted in MongoDB
- ✅ JWT tokens validated on every request
- ✅ Auth0 audience verification
- ✅ MongoDB over TLS/SSL
- ✅ CORS configured
- ✅ Secrets in environment variables

---

## 📞 Support Files

- **MongoDB Setup:** [`MONGODB_SETUP.md`](MONGODB_SETUP.md)
- **Full Integration Guide:** [`INTEGRATION_COMPLETE.md`](INTEGRATION_COMPLETE.md)
- **Auth0 Setup:** [`AUTH0_SETUP.md`](../chatapp/AUTH0_SETUP.md)

---

**You're all set! 🎉** Messages are now encrypted and stored in MongoDB!
