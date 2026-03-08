# Clear MongoDB collections for fresh start
Write-Host "Connecting to MongoDB..." -ForegroundColor Cyan

# Find mongosh or mongo executable
$mongosh = Get-Command mongosh -ErrorAction SilentlyContinue
$mongo = Get-Command mongo -ErrorAction SilentlyContinue

$mongoCmd = $null
if ($mongosh) {
    $mongoCmd = "mongosh"
} elseif ($mongo) {
    $mongoCmd = "mongo"
} else {
    Write-Host "Error: MongoDB shell (mongosh or mongo) not found in PATH" -ForegroundColor Red
    Write-Host "Please install MongoDB shell or add it to your PATH" -ForegroundColor Yellow
    exit 1
}

Write-Host "Using $mongoCmd" -ForegroundColor Green

# Clear conversations and messages
& $mongoCmd "mongodb://localhost:27017/chatapp_db" --quiet --eval @"
db.conversations.deleteMany({});
db.messages.deleteMany({});
print('✅ Cleared all conversations and messages');
print('Collections cleared:');
print('- Conversations: ' + db.conversations.countDocuments({}));
print('- Messages: ' + db.messages.countDocuments({}));
"@

Write-Host "`n✅ Database cleared! Now restart your Flutter app and send a new message." -ForegroundColor Green
