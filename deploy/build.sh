#!/bin/bash
set -e

echo "Building BlueJay CMS..."
go build -o bluejay-cms cmd/server/main.go

echo "Build complete!"
