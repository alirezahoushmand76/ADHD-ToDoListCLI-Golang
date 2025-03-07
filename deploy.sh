#!/bin/bash

# This script deploys the todolist application on a server
# It pulls the latest Docker image and runs it

# Configuration
DOCKER_IMAGE="yourusername/todolist:latest"  # Replace with your Docker Hub username
VOLUME_NAME="todolist_data"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed"
    exit 1
fi

# Pull the latest image
echo "Pulling the latest Docker image..."
docker pull $DOCKER_IMAGE

# Check if the volume exists, create it if it doesn't
if ! docker volume inspect $VOLUME_NAME &> /dev/null; then
    echo "Creating Docker volume $VOLUME_NAME..."
    docker volume create $VOLUME_NAME
fi

# Create a wrapper script
echo "Creating wrapper script..."
cat > todolist << EOF
#!/bin/bash
docker run --rm -v $VOLUME_NAME:/home/appuser/.todolist $DOCKER_IMAGE "\$@"
EOF

# Make the wrapper script executable
chmod +x todolist

# Move the wrapper script to a location in PATH
echo "Installing the wrapper script..."
sudo mv todolist /usr/local/bin/

echo "Deployment completed successfully!"
echo "You can now use the 'todolist' command to run the application." 