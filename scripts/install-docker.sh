#!/bin/bash
# Docker installation for Debian with cleanup
# Usage: ./install-docker.sh

set -e  # Exit on any error

echo "Checking Docker installation..."

# Check if Docker is already installed and working
if command -v docker &> /dev/null; then
    echo "Docker is already installed. Version:"
    docker --version
    
    # Test if it's working
    if sudo docker ps &> /dev/null; then
        echo "Docker is working correctly."
        exit 0
    fi
fi

echo "Installing/fixing Docker installation..."

# Step 1: Remove any existing Docker installations and repositories
echo "Cleaning up any existing Docker installations..."
sudo apt-get remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true

# Remove any existing Docker repositories that might be misconfigured
echo "Removing existing Docker repositories..."
sudo rm -f /etc/apt/sources.list.d/docker.list
sudo rm -f /etc/apt/keyrings/docker.gpg
sudo rm -f /usr/share/keyrings/docker-archive-keyring.gpg

# Step 2: Update package lists (this should work now without the bad repos)
echo "Updating package lists..."
sudo apt-get update

# Step 3: Install prerequisites
echo "Installing prerequisites..."
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# Step 4: Add Docker's official GPG key for Debian
echo "Adding Docker GPG key..."
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Step 5: Get the correct Debian codename
DEBIAN_CODENAME=$(lsb_release -cs)
echo "Detected Debian codename: $DEBIAN_CODENAME"

# Handle special case where lsb_release might return something Docker doesn't recognize
case $DEBIAN_CODENAME in
    "bookworm")
        DOCKER_CODENAME="bookworm"
        ;;
    "bullseye")
        DOCKER_CODENAME="bullseye"
        ;;
    "buster")
        DOCKER_CODENAME="buster"
        ;;
    *)
        echo "Warning: Unrecognized Debian version, using bullseye as fallback"
        DOCKER_CODENAME="bullseye"
        ;;
esac

# Step 6: Add Docker repository for Debian
echo "Adding Docker repository for Debian $DOCKER_CODENAME..."
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  $DOCKER_CODENAME stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Step 7: Update package lists again
echo "Updating package lists with new Docker repository..."
sudo apt-get update

# Step 8: Install Docker
echo "Installing Docker..."
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Step 9: Add user to docker group
echo "Adding user $USER to docker group..."
sudo usermod -aG docker $USER

# Step 10: Start and enable Docker
echo "Starting Docker service..."
sudo systemctl start docker
sudo systemctl enable docker

echo "Docker installation completed successfully!"
echo "Docker version:"
sudo docker --version

# Test Docker installation
echo "Testing Docker installation..."
if sudo docker run --rm hello-world &> /dev/null; then
    echo "Docker is working correctly!"
    echo ""
    echo "IMPORTANT: Log out and log back in (or run 'newgrp docker') to use Docker without sudo."
else
    echo "Docker installation may have issues, but basic installation completed."
fi

echo "Installation complete!"