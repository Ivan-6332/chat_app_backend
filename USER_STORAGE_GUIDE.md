# User Storage & Message Association - Implementation Guide

## 🎉 What Was Implemented

Your chat app now **automatically stores users in MongoDB** when they log in and send messages. Every message is associated with the user who sent it!

---

## 🔄 How It Works

### 1. **User Login Flow**

When a user logs in via Auth0:
1. Flutter app authenticates with Auth0
2. Receives JWT token with user claims (sub, email, name)
3. Token is stored securely in Flutter

### 2. **Sending a Message Flow**

When a user sends a message:
1. **Flutter App:**
   - Encrypts the message locally
   - Sends encrypted message with JWT token to backend
   
2. **Backend API:**
   - Validates JWT token (Auth0 middleware)
   - Extracts user info from JWT claims:
     - `sub` → Auth0 ID (unique identifier)
     - `email` → User's email
     - `name` or `nickname` → Username
   - **Checks MongoDB `users` collection:**
     - If user exists → Updates last_seen timestamp
     - If user doesn't exist → Creates new user document
   - Saves encrypted message to `messages` collection
   - Links message to user via `sender_id` field
   - Returns response with message and user data

### 3. **Data Structure**

#### Users Collection
```json
{
  "_id": ObjectId("..."),
  "auth0_id": "auth0|123456789",
  "username": "john_doe",
  "email": "john@example.com",
  "display_name": "",
  "profile_pic": "",
  "public_key": "",
  "created_at": ISODate("2026-03-08T..."),
  "last_seen_at": ISODate("2026-03-08T...")
}
```

#### Messages Collection
```json
{
  "_id": ObjectId("..."),
  "sender_id": "auth0|123456789",  // Links to user
  "conversation_id": "conv_123",
  "encrypted_text": "U2FsdGVkX1...",
  "timestamp": ISODate("2026-03-08T...")
}
```

#### Conversations Collection
```json
{
  "_id": ObjectId("..."),
  "members": ["auth0|123", "auth0|456"],
  "last_message": "U2FsdGVkX1...",
  "last_message_at": ISODate("2026-03-08T..."),
  "created_at": ISODate("2026-03-08T...")
}
```

---

## 📡 New API Endpoints

### User Endpoints

#### Get Current User
```http
GET /api/v1/users/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "65f1a2b3c4d5e6f7g8h9i0j1",
    "auth0Id": "auth0|123456789",
    "username": "john_doe",
    "email": "john@example.com",
    "createdAt": "2026-03-08T10:30:00Z",
    "lastSeenAt": "2026-03-08T14:25:00Z"
  }
}
```

#### Get User by ID
```http
GET /api/v1/users/:id
Authorization: Bearer <token>
```

### Message Endpoint (Updated)

#### Send Message
```http
POST /api/v1/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "senderId": "auth0|123456789",
  "conversationId": "conv_123",
  "encryptedText": "U2FsdGVkX1..."
}
```

**Response (Now includes user info):**
```json
{
  "success": true,
  "data": {
    "message": {
      "id": "65f1a2b3c4d5e6f7g8h9i0j1",
      "senderId": "auth0|123456789",
      "conversationId": "conv_123",
      "encryptedText": "U2FsdGVkX1...",
      "timestamp": "2026-03-08T14:25:00Z"
    },
    "user": {
      "id": "65f1a2b3c4d5e6f7g8h9i0j1",
      "auth0Id": "auth0|123456789",
      "username": "john_doe",
      "email": "john@example.com"
    }
  }
}
```

---

## 💻 Flutter Usage

### Get Current User from Backend
```dart
import 'services/chat_service.dart';

// After login, get user profile
final user = await ChatService.getCurrentUser();
print('Logged in as: ${user?['username']}');
```

### Send Message (Automatic User Creation)
```dart
// User will be automatically created/updated when sending
await ChatService.sendMessage(
  plainText: 'Hello!',
  conversationId: 'conv_123',
  senderId: 'auth0|123456789',
);
```

### Get Another User's Profile
```dart
final otherUser = await ChatService.getUserById('auth0|987654321');
print('User: ${otherUser?['username']}');
```

---

## 🔍 View Your Data in MongoDB

1. Go to [MongoDB Atlas Dashboard](https://cloud.mongodb.com)
2. Click **"Browse Collections"**
3. Select `chatapp_db` database
4. View collections:

### Users Collection
- See all registered users
- Check when they were created
- See last login time

### Messages Collection
- All messages stored encrypted
- Each message has `sender_id` linking to user
- Sorted by timestamp

### Conversations Collection
- Chat conversations
- Members list shows which users are in each chat
- Last message preview

---

## 🔐 Security Features

✅ **JWT Validation** - Every request validates Auth0 token  
✅ **Automatic User Creation** - No manual registration needed  
✅ **Encrypted Storage** - Messages stored encrypted in database  
✅ **User Isolation** - Each message linked to authenticated user  
✅ **No Password Storage** - Auth0 handles authentication  
✅ **Secure Claims** - User ID comes from validated JWT token  

---

## 🧪 Testing Your Implementation

### Test 1: User Creation on First Message
1. Login with a new Auth0 account
2. Send a message
3. Check MongoDB users collection
4. You should see a new user document

### Test 2: User Update on Subsequent Messages
1. Send another message
2. Check the user's `last_seen_at` field
3. It should be updated to current time

### Test 3: Retrieve User Profile
```dart
// In your Flutter app
final user = await ChatService.getCurrentUser();
if (user != null) {
  debugPrint('User ID: ${user['id']}');
  debugPrint('Username: ${user['username']}');
  debugPrint('Email: ${user['email']}');
}
```

### Test 4: View Associated Messages
```sql
// In MongoDB Compass or Atlas
db.messages.find({ "sender_id": "auth0|YOUR_USER_ID" })
```

---

## 📊 Database Query Examples

### Find all messages by a user
```javascript
db.messages.find({ sender_id: "auth0|123456789" })
```

### Find user by email
```javascript
db.users.findOne({ email: "john@example.com" })
```

### Find recent messages in a conversation
```javascript
db.messages.find({ 
  conversation_id: "conv_123" 
}).sort({ timestamp: -1 }).limit(50)
```

### Get all users who sent messages today
```javascript
db.users.find({ 
  last_seen_at: { 
    $gte: new Date("2026-03-08T00:00:00Z") 
  }
})
```

---

## 🎯 What Happens Behind the Scenes

### When User Logs In:
1. Auth0 creates JWT token with user claims
2. Flutter stores token securely
3. **No immediate database action** (user created on first message)

### When User Sends First Message:
1. Backend validates JWT token
2. Checks if user exists in `users` collection
3. **User doesn't exist** → Creates new user from JWT claims
4. Saves message with `sender_id`
5. Returns message + user info

### When User Sends Another Message:
1. Backend validates JWT token
2. Checks if user exists in `users` collection
3. **User exists** → Updates `last_seen_at` timestamp
4. Saves message
5. Returns message + user info

---

## 🚀 Benefits of This Approach

✅ **Automatic** - No manual user registration API calls  
✅ **Seamless** - Works transparently with Auth0  
✅ **Efficient** - User created only when needed  
✅ **Secure** - User data from validated JWT tokens  
✅ **Trackable** - Know when users last used the app  
✅ **Relational** - Messages linked to users properly  

---

## 🔧 Configuration Files

All the necessary files have been updated:

### Backend:
- ✅ `controllers/message_controller.go` - Handles user creation on message send
- ✅ `controllers/user_controller.go` - New user endpoints
- ✅ `services/user_service.go` - User database operations
- ✅ `models/user.go` - User data model
- ✅ `routes/routes.go` - User API routes
- ✅ `main.go` - Wired up user service

### Flutter:
- ✅ `lib/services/api_service.dart` - User API calls
- ✅ `lib/services/chat_service.dart` - User retrieval methods

---

## 📝 Summary

**Now when a user:**
1. Logs in → JWT token stored
2. Sends message → User automatically created in MongoDB
3. Sends more messages → User's last_seen updated
4. All messages linked to that user via sender_id

**You can:**
- See all users in MongoDB `users` collection
- Track when users were created and last seen
- Query messages by user
- Get user profiles via API
- Everything happens automatically!

---

## ✨ Next Steps

Consider adding:
- User profile pictures
- User status (online/offline)
- User bio/description
- Friend lists
- Block/unblock users
- User search functionality

**Your chat app now has full user management with MongoDB! 🎉**
