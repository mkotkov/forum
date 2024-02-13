#!/bin/bash

# Define variables
DOCKER_IMAGE_NAME="forum"
DOCKER_IMAGE_TAG="latest"
DOCKER_FILE="Dockerfile"

# Build Docker image
docker build -t $DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG -f $DOCKER_FILE .

# Run Docker container
docker run -d -p 8080:8080 $DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG
