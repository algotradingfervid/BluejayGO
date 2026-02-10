#!/bin/bash
# Kill process on port 28090, rebuild, and restart

echo "Killing process on port 28090..."
lsof -ti:28090 | xargs kill -9 2>/dev/null
echo "Building..."
cd "$(dirname "$0")"
go build -o server ./cmd/server
if [ $? -eq 0 ]; then
    echo "Starting server on port 28090..."
    ./server
else
    echo "Build failed!"
    exit 1
fi
