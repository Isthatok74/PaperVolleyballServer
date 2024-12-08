@echo off

rem Change to this script's directory
cd /d "%~dp0"

rem Generate Go code from `server.proto`
protoc --go_out=go --go-grpc_out=go server.proto

rem Generate C# code from `server.proto`
protoc --csharp_out=cs server.proto