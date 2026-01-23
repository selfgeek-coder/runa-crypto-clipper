@echo off
chcp 65001 >nul

setlocal enabledelayedexpansion

rem go check

where go >nul 2>&1
if %errorlevel%==0 (
    rem installed
) else (
    echo Go not found. Install go and add it to PATH / Go не найден. Установите Go и добавьте его в PATH.
    echo Visit https://go.dev/dl/ for download / Перейдите на https://go.dev/dl/ для скачивания
    pause
    exit /b 1
)

rem garble check

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

rem upx check

if exist upx.exe (
    rem installed
) else (
    echo Downloading UPX / Загрузка UPX...

    powershell -NoProfile -Command ^
        "Invoke-WebRequest -Uri 'https://github.com/upx/upx/releases/download/v5.1.0/upx-5.1.0-win64.zip' -OutFile 'upx.zip'" || (
        echo Failed to download UPX / Ошибка при скачивании UPX
        pause
        exit /b 1
    )

    powershell -NoProfile -Command ^
        "Expand-Archive -Force 'upx.zip' 'upx_tmp'" || (
        echo Failed to unzip UPX / Ошибка при распаковке UPX
        pause
        exit /b 1
    )

    for /r upx_tmp %%f in (upx.exe) do (
        copy "%%f" ".\upx.exe" >nul
    )

    rd /s /q upx_tmp
    del /q upx.zip
)

pause & cls