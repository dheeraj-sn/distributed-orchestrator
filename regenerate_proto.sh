#!/bin/bash

# Exit on any error
set -e

# Directory containing this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Ensure we're in the right directory
cd "$DIR"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please install Protocol Buffers compiler first"
    exit 1
fi

# Check if Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null || ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Error: protoc-gen-go and/or protoc-gen-go-grpc not found"
    echo "Please install them with:"
    echo "go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    echo "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

echo "Generating protobuf code..."

# Generate the protobuf and gRPC code
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    proto/orchestrator.proto

echo "Protocol buffer code generation completed successfully"