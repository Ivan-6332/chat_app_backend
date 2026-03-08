# User Storage - Quick Reference

## ✅ What Changed

### Automatic User Storage
- When user sends a message → User is **automatically** created in MongoDB `users` collection
- User data is extracted from Auth0 JWT token
- No manual registration needed!

### Collections Structure

**users** → Stores user profiles
```
auth0_id | username | email | created_at | last_seen_at
```

**messages** → Linked to users via sender_id
```
sender_id | conversation_id | encrypted_text | timestamp
```

**conversations** → Members array contains user auth0_ids
```
members[] | last_message | last_message_at | created_at
```

---

## 🚀 Test It Now

### 1. View MongoDB Collections
```
1. Go to MongoDB Atlas → Browse Collections
2. Database: chatapp_db
3. Collections: users, messages, conversations
```

### 2. Send a Test Message
```dart
// Flutter app
await ChatService.sendMessage(
  plainText: 'Test message',
  conversationId: 'test_conv',
  senderId: 'auto_from_token',  // Auto-replaced by backend
);
```

### 3. Check User Was Created
MongoDB Query:
```javascript
db.users.find()
```

You should see:
```json
{
  "auth0_id": "auth0|...",
  "username": "your_name",
  "email": "your@email.com",
  "created_at": "2026-03-08...",
  "last_seen_at": "2026-03-08..."
}
```

### 4. Check Message Was Saved
```javascript
db.messages.find({ sender_id: "your_auth0_id" })
```

---

## 📡 New Endpoints

```
GET  /api/v1/users/me        → Get current logged-in user
GET  /api/v1/users/:id       → Get user by Auth0 ID
POST /api/v1/messages        → Send message (creates user automatically)
GET  /api/v1/messages/:convId → Get messages
GET  /api/v1/conversations/:userId → Get conversations
```

---

## 🔥 Key Features

✅ Users **auto-created** from Auth0 JWT on first message  
✅ Last seen timestamp **auto-updated** on each message  
✅ Messages **linked to users** via sender_id  
✅ All data **encrypted** in database  
✅ Full **user profiles** available via API  

---

## 💡 Usage Examples

### Flutter: Get Current User
```dart
final user = await ChatService.getCurrentUser();
print('${user?['username']}'); // "john_doe"
```

### Flutter: Send Message
```dart
// User automatically created/updated
await ChatService.sendMessage(
  plainText: 'Hello World',
  conversationId: 'conv_123',
  senderId: userId,
);
```

### Backend Response
```json
{
  "message": { /* message data */ },
  "user": {
    "auth0Id": "auth0|123",
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

---

## 🔍 Flow Diagram

```
User Logs In (Auth0)
     ↓
JWT Token Generated
     ↓
User Sends Message
     ↓
Backend Receives Request
     ↓
JWT Token Validated ✅
     ↓
Extract User Info from Token
     ↓
Check MongoDB Users Collection
     ↙         ↘
User Exists    User Doesn't Exist
     ↓              ↓
Update         Create User
last_seen      in Database
     ↓              ↓
     └──────┬───────┘
            ↓
    Save Message
    (with sender_id)
            ↓
    Return Response
```

---

## 🎯 Complete Setup Checklist

✅ MongoDB connected  
✅ Auth0 configured with audience  
✅ User model created  
✅ User service implemented  
✅ User controller created  
✅ Message controller updated  
✅ Routes added for users  
✅ Flutter API service updated  
✅ Flutter chat service updated  
✅ Backend running on port 8080  

**Everything is ready to test! 🚀**

---

## 📝 Quick MongoDB Queries

```javascript
// Find all users
db.users.find()

// Find messages by user
db.messages.find({ sender_id: "auth0|123" })

// Find user's conversations
db.conversations.find({ members: "auth0|123" })

// Count total users
db.users.countDocuments()

// Find users active today
db.users.find({ 
  last_seen_at: { $gte: new Date("2026-03-08") } 
})
```

---

**That's it! Your users are now automatically stored in MongoDB! 🎉**
