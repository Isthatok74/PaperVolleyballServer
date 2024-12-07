@echo off
rem Set the directory for binaries
set BINDIR=build

rem Ensure the build directory exists
mkdir %BINDIR% 2>nul

rem Run the Go build command
go build -o %BINDIR% ./cmd/pv-server