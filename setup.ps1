# Setup Script for ChatApp Backend

Write-Host "🚀 Setting up ChatApp Backend..." -ForegroundColor Cyan

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go is not installed. Please install Go 1.21 or higher." -ForegroundColor Red
    exit 1
}

Write-Host "✅ Go version:" (go version) -ForegroundColor Green

# Check if .env file exists
if (-not (Test-Path ".env")) {
    Write-Host "⚠️  .env file not found. Creating from .env.example..." -ForegroundColor Yellow
    Copy-Item ".env.example" ".env"
    Write-Host "📝 Please edit .env file with your Auth0 credentials" -ForegroundColor Yellow
} else {
    Write-Host "✅ .env file found" -ForegroundColor Green
}

# Install dependencies
Write-Host "`n📦 Installing dependencies..." -ForegroundColor Cyan
go mod download
go mod tidy

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Dependencies installed successfully" -ForegroundColor Green
} else {
    Write-Host "❌ Failed to install dependencies" -ForegroundColor Red
    exit 1
}

# Build the application
Write-Host "`n🔨 Building application..." -ForegroundColor Cyan
go build -o chatapp-backend.exe .

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build successful!" -ForegroundColor Green
    Write-Host "`n✨ Setup complete!" -ForegroundColor Green
    Write-Host "`nTo start the server, run:" -ForegroundColor Cyan
    Write-Host "  .\chatapp-backend.exe" -ForegroundColor White
    Write-Host "`nOr use:" -ForegroundColor Cyan
    Write-Host "  go run main.go" -ForegroundColor White
} else {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
}
