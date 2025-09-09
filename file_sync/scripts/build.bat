@echo off
echo Building Windows executable...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-H windowsgui" -o file-sync.exe
pause