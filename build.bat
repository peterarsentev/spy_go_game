#!/bin/bash

set -e

echo "Building app for Linux..."
set GOOS=linux
set GOARCH=amd64
go build -o spy .\cmd\main.go
