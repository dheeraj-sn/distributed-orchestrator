#!/bin/bash

# Start the scheduler
go run cmd/scheduler/main.go &
SCHEDULER_PID=$!

# Wait a bit for scheduler to start
sleep 2

# Start a worker
go run cmd/worker/main.go &
WORKER_PID=$!

# Handle cleanup on script exit
cleanup() {
    echo "Shutting down services..."
    kill $SCHEDULER_PID
    kill $WORKER_PID
    exit 0
}

trap cleanup SIGINT SIGTERM

# Keep script running
wait 