#!/bin/bash

echo "🚀 Setting up ChatApp Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "⚠️  .env file not found. Creating from .env.example..."
    cp .env.example .env
    echo "📝 Please edit .env file with your Auth0 credentials"
else
    echo "✅ .env file found"
fi

# Install dependencies
echo ""
echo "📦 Installing dependencies..."
go mod download
go mod tidy

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

# Build the application
echo ""
echo "🔨 Building application..."
go build -o chatapp-backend .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "✨ Setup complete!"
    echo ""
    echo "To start the server, run:"
    echo "  ./chatapp-backend"
    echo ""
    echo "Or use:"
    echo "  go run main.go"
else
    echo "❌ Build failed"
    exit 1
fi
