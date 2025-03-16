#!/bin/bash

# Find and kill the process
PID=$(pgrep -f "go-image-server")
if [ -z "$PID" ]; then
    echo "go-image-server is not running"
    exit 1
fi

echo "Stopping go-image-server (PID: $PID)..."
kill $PID

# Wait and verify
sleep 2
if pgrep -f "go-image-server" > /dev/null; then
    echo "Force stopping go-image-server..."
    pkill -9 -f "go-image-server"
else
    echo "go-image-server stopped successfully"
fi
