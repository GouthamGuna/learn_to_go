# build without console window (for service)

Write-Host "Building Windows executable..." -ForegroundColor Green

# Set environment variables
$env:GOOS = "windows"
$env:GOARCH = "amd64" 
$env:CGO_ENABLED = "0"

# Build the executable
go build -ldflags="-H windowsgui" -o file-sync.exe

# Check if build was successful
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful! Created file-sync.exe" -ForegroundColor Green
} else {
    Write-Host "Build failed!" -ForegroundColor Red
}

# Pause to see the output
Read-Host "Press Enter to continue"