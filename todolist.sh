#!/bin/bash

# This script is a wrapper for the dockerized todolist application
# It allows you to run the application as if it were installed locally

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Error: Docker Compose is not installed or not in PATH"
    exit 1
fi

# Check if the image exists, if not build it
if ! docker image inspect todolist:latest &> /dev/null; then
    echo "Building todolist Docker image..."
    docker-compose build
fi

# Run the command
docker-compose run --rm todolist "$@" 