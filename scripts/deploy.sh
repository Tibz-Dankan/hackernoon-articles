#!/bin/bash
# Zero-downtime deployment script using Docker Swarm
# Usage: ./deploy.sh [new-image-tag]

# Set the app environment variable
export APP_ENV="development"

set -e  # Exit on any error

NEW_IMAGE=$1
if [ -z "$NEW_IMAGE" ]; then
  echo "Error: New image tag not provided"
  echo "Usage: ./deploy.sh [new-image-tag]"
  exit 1
fi

echo "Deploying new image: $NEW_IMAGE"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "ERROR: Docker is not installed!"
    echo "Please run ./install-docker.sh first to install Docker"
    exit 1
fi

# Check if current user can run Docker commands
if ! docker info &> /dev/null; then
    echo "ERROR: Cannot run Docker commands. This could be because:"
    echo "1. Docker service is not running"
    echo "2. Current user is not in the docker group"
    echo "3. You need to log out and back in after adding user to docker group"
    echo ""
    echo "Try running: sudo systemctl start docker"
    echo "Or run: newgrp docker"
    exit 1
fi

# Create the project directory if it doesn't exist
mkdir -p ./app
cd app

# Copy .env file from the root directory to app directory
if [ -f ~/.env ]; then
    echo "Copying .env file from home directory to app directory..."
    cp ~/.env ./
else
    echo "ERROR: No .env file found in .app/ !"
    echo "Please create an .env file with your application's environment variables"
    exit 1
fi

# Copy docker-compose.yaml file from the root directory to app directory
if [ -f ~/docker-compose.yaml ]; then
    echo "Copying docker-compose.yaml from home directory to app directory..."
    cp ~/docker-compose.yaml ./
else
    echo "ERROR: No docker-compose.yaml file found in ~/ !"
    echo "Please create an docker-compose.yaml file with deployment configuration"
    exit 1
fi

# Initialize Docker Swarm if not already initialized
if ! docker info | grep -q "Swarm: active"; then
  echo "Initializing Docker Swarm..."
  docker swarm init --advertise-addr $(hostname -I | awk '{print $1}')
fi

echo "About to deploy stack with image: ${NEW_IMAGE}"

# Export NEW_IMAGE as environment variable for docker-compose
export NEW_IMAGE=$NEW_IMAGE

# Deploy or update the stack
if docker stack ls | grep -q "app-stack"; then
  echo "Updating existing stack..."
  docker stack deploy -c docker-compose.yaml app-stack --with-registry-auth
else
  echo "Deploying new stack..."
  docker stack deploy -c docker-compose.yaml app-stack --with-registry-auth
fi

# Wait for backend service to be running and available on port 3000
echo "Waiting for backend service to start and be available on port 3000..."

# # Function to check if service is running properly
# check_service() {
#   # Check if service exists and is running with correct replica count
#   if docker service ls | grep app-stack_hackernoon-index | grep -q "1/1"; then
#     # Check if port 3000 is listening
#     if timeout 5 bash -c "</dev/tcp/localhost/3000" &>/dev/null; then
#       # Check if health endpoint returns 200 OK
#       if curl -s -f -o /dev/null http://localhost:3000/health; then
#         return 0  # Service is running properly
#       fi
#     fi
#   fi
#   return 1  # Service is not running properly
# }

check_service() {
  echo "  Checking service existence..."
  if ! docker service ls | grep -q app-stack_hackernoon-index; then
    echo "  Service app-stack_hackernoon-index not found"
    return 1
  fi
  
  echo "  Checking service replicas..."
  if ! docker service ls | grep app-stack_hackernoon-index | grep -q "1/1"; then
    echo "  Service replicas not ready yet"
    return 1
  fi
  
  echo "  Checking port 3000..."
  if ! timeout 5 bash -c "</dev/tcp/localhost/3000" &>/dev/null; then
    echo "  Port 3000 not accessible"
    return 1
  fi
  
  echo "  Checking health endpoint..."
  if ! curl -s -f -o /dev/null http://localhost:3000/health; then
    echo "  Health endpoint not responding"
    return 1
  fi
  
  return 0
}

# Keep checking until service is running or timeout 
MAX_ATTEMPTS=40
ATTEMPT=1
WAIT_TIME=5 # seconds between attempts

while [ $ATTEMPT -le $MAX_ATTEMPTS ]; do
  echo "Checking if backend service is running (attempt $ATTEMPT/$MAX_ATTEMPTS)..."
  
  if check_service; then
    echo "Backend service is now running and available on port 3000!"
    break
  fi
  
  # If this is the first few attempts, show more debugging info
  if [ $ATTEMPT -le 3 ] || [ $(($ATTEMPT % 10)) -eq 0 ]; then
    echo "Service status:"
    docker service ps app-stack_hackernoon-index --no-trunc
    echo "Recent logs:"
    docker service logs app-stack_hackernoon-index --tail 10
  fi
  
  ATTEMPT=$((ATTEMPT+1))
  
  if [ $ATTEMPT -gt $MAX_ATTEMPTS ]; then
    echo "ERROR: Backend service failed to start within the allocated time."
    echo "Final service status:"
    docker service ps app-stack_hackernoon-index --no-trunc
    echo "Detailed container logs:"
    docker service logs app-stack_hackernoon-index --tail 100
    echo "Checking if .env file variables are being properly loaded..."
    echo "Number of variables in .env file: $(grep -v '^#' ./.env | grep -v '^$' | wc -l)"
    exit 1
  fi
  
  echo "Waiting for $WAIT_TIME seconds before next check..."
  sleep $WAIT_TIME
done

echo "Development Deployment successful!"
echo "Traefik dashboard is available on port 8080 (username: admin, password: admin123)"