# MongoDB Atlas Setup Guide

## Step 1: Create MongoDB Atlas Account

1. Go to [MongoDB Atlas](https://www.mongodb.com/cloud/atlas/register)
2. Sign up for a **FREE** account (no credit card required)
3. Verify your email address

## Step 2: Create a New Cluster

1. After logging in, click **"Build a Database"**
2. Choose **"M0 FREE"** tier (shared cluster)
   - Provider: **AWS** (or your preferred cloud provider)
   - Region: Choose closest to your location (e.g., **us-east-1**)
3. Cluster Name: `ChatAppCluster` (or any name you prefer)
4. Click **"Create"** (this takes 1-3 minutes)

## Step 3: Create Database User

1. Click **"Database Access"** in the left sidebar (under Security)
2. Click **"+ Add New Database User"**
3. Authentication Method: **Password**
   - Username: `chatapp_user` (or your choice)
   - Password: Click **"Autogenerate Secure Password"** and **SAVE IT SOMEWHERE SAFE**
   - Or create your own strong password
4. Database User Privileges: **Read and write to any database**
5. Click **"Add User"**

**⚠️ IMPORTANT: Copy and save your password now! You'll need it for the connection string.**

## Step 4: Configure Network Access

1. Click **"Network Access"** in the left sidebar (under Security)
2. Click **"+ Add IP Address"**
3. For development, click **"Allow Access from Anywhere"** (0.0.0.0/0)
   - ⚠️ For production, restrict to specific IPs
4. Click **"Confirm"**

## Step 5: Get Connection String

1. Click **"Database"** in the left sidebar
2. Click **"Connect"** button on your cluster
3. Choose **"Connect your application"**
4. Driver: **Go** | Version: **1.17 or later**
5. Copy the connection string, it looks like:
   ```
   mongodb+srv://chatapp_user:<password>@chatappcluster.xxxxx.mongodb.net/?retryWrites=true&w=majority
   ```
6. **Replace `<password>` with your actual database user password**

## Step 6: Create Database and Collections (Auto-created by app)

The application will automatically create:
- Database: `chatapp_db`
- Collections: `users`, `conversations`, `messages`

## Step 7: Add Connection String to Your Backend

Add to your `.env` file:
```env
MONGODB_URI=mongodb+srv://chatapp_user:YOUR_PASSWORD@chatappcluster.xxxxx.mongodb.net/?retryWrites=true&w=majority
MONGODB_DATABASE=chatapp_db
```

**Replace `YOUR_PASSWORD` with your actual password!**

## Verify Connection

Once configured, run your backend and check the logs for:
```
✅ Connected to MongoDB successfully
```

## Optional: Browse Your Data

1. In Atlas Dashboard, click **"Browse Collections"**
2. You'll see your database and collections
3. You can view, edit, and query data directly from the web interface

---

## Important Security Notes

🔒 **For Production:**
- Use IP whitelisting instead of "Allow Access from Anywhere"
- Use environment variables for credentials (never hardcode)
- Enable MongoDB authentication
- Use strong passwords
- Regularly rotate credentials
