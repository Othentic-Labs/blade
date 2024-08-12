#!/usr/bin/env bash

# Function to show help
show_help() {
    echo "Usage: run.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --destroy    Destroy the environment"
    echo "  --local      Run docker-compose with local configuration"
    echo "  --holesky    Run docker-compose with Holesky configuration"
    echo "  --help       Display this help message"
}

# Function to destroy the environment
destroy_environment() {
    echo "Destroying the environment..."

    cd docker/local
    
    # Stop and remove all containers defined in the docker-compose file
    docker-compose down
    
    # Remove all images defined in the docker-compose file
    docker-compose rm -f -v
    
    # Remove all unused images
    docker image prune -af
    
    # Remove all unused volumes
    docker volume prune -f
    
    # Remove specific volumes if needed
    docker volume rm $(docker volume ls -qf dangling=true)
    
    # Optionally, remove network if it exists and is not needed
    docker network rm blade-docker-network
    
    # Optionally, remove local files if needed (adjust paths as necessary)
    # For example, remove data directories
    rm -rf docker/local/services/blockscout-db-data
    rm -rf docker/local/services/logs
    rm -rf docker/local/services/redis-data
    rm -rf docker/local/services/stats-db-data
    
    echo "Docker environment destroyed successfully."
}

# Parse command line arguments
case "$1" in
    --destroy)
        destroy_environment
        ;;
    --local)
        echo "Running local environment..."
        cd docker/local
        docker-compose up -d
        ;;
    --holesky)
        echo "Running Holesky environment..."
        cd docker/local
        docker-compose -f docker-compose_holesky.yml up -d
        ;;
    --help|*)
        show_help
        ;;
esac
