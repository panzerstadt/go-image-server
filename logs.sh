#!/bin/bash

# Check if logs directory exists
if [ ! -d "logs" ]; then
    echo "Logs directory doesn't exist. Has the server been started?"
    exit 1
fi

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo "tmux is not installed. Using simple log view..."
    echo "Press Ctrl+C to exit"
    echo "=== Output Log ==="
    tail -f logs/output.log &
    echo "=== Error Log ==="
    tail -f logs/error.log
    exit 0
fi


# Check if tmux session exists and attach to it if it does
if tmux has-session -t go-image-logs 2>/dev/null; then
    # Session exists, just attach to it
    tmux attach-session -t go-image-logs
else
    # Create a new tmux session if not already in one
    if [ -z "$TMUX" ]; then
        tmux new-session -d -s go-image-logs
        tmux split-window -h
        tmux select-pane -t 0
        tmux send-keys "echo '=== Output Log ===' && tail -f logs/output.log" C-m
        tmux select-pane -t 1
        tmux send-keys "echo '=== Error Log ===' && tail -f logs/error.log" C-m
        tmux attach-session -t go-image-logs
    else
        # Already in tmux, just create splits
        tmux split-window -h
        tmux select-pane -t 0
        tmux send-keys "echo '=== Output Log ===' && tail -f logs/output.log" C-m
        tmux select-pane -t 1
        tmux send-keys "echo '=== Error Log ===' && tail -f logs/error.log" C-m
    fi
fi
