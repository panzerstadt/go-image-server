#!/bin/bash

# Create logs directory if it doesn't exist
mkdir -p logs

# Check if server is already running
if pgrep -f "go-image-server" > /dev/null; then
    echo "go-image-server is already running!"
    exit 1
fi

# Start the server
nohup ./go-image-server > logs/output.log 2> logs/error.log &

# Get the PID of the new process
PID=$!
echo "go-image-server started with PID: $PID"
echo "Logs are available at:"
echo "  - logs/output.log (stdout)"
echo "  - logs/error.log (stderr)"
