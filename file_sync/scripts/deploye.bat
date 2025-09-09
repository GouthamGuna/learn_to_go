@echo off
echo Installing and starting File Sync Service...

:: Reference the parent directory of the script for file-sync.exe
SET EXE_PATH=%~dp0..\file-sync.exe

IF EXIST "%EXE_PATH%" (
    "%EXE_PATH%" install
    "%EXE_PATH%" start
) ELSE (
    echo file-sync.exe not found in %EXE_PATH%
    echo Failed to install service!
    exit /b 1
)

echo Service installed and started successfully!
echo.
pause