#!/bin/bash

set -e

echo "Building app for Linux..."
set GOOS=linux
set GOARCH=amd64
go build -o spy.1.0.1 .\cmd\main.go
