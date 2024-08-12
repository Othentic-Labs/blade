#!/usr/bin/env bash

# Function to check if a command is available
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check dependencies
check_dependencies() {
    local missing_deps=()

    # Check for Docker
    if ! command_exists docker; then
        missing_deps+=("Docker")
    fi

    # Check for Docker Compose
    if ! command_exists docker-compose; then
        missing_deps+=("Docker Compose")
    fi

    # Add checks for any other dependencies here
    # For example:
    # if ! command_exists jq; then
    #     missing_deps+=("jq")
    # fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo "Error: The following dependencies are missing:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo "Please install the missing dependencies and try again."
        exit 1
    fi

    echo "All dependencies are installed."
}

show_help() {
    echo "Usage: run.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --destroy    Destroy the environment"
    echo "  --local      Run blade network with local configuration"
    echo "  --holesky    Run blade network with Holesky L1"
    echo "  --help       Display this help message"
}

destroy_environment() {
    echo "Destroying the environment..."

    cd docker/local
    
    docker-compose down
    
    docker-compose rm -f -v
    
    docker image prune -af
    
    docker volume prune -f
    
    docker volume rm $(docker volume ls -qf dangling=true)
    
    docker network rm blade-docker-network
    
    rm -rf docker/local/services/blockscout-db-data
    rm -rf docker/local/services/logs
    rm -rf docker/local/services/redis-data
    rm -rf docker/local/services/stats-db-data
    
    echo "Docker environment destroyed successfully."
}

# Check dependencies before proceeding
check_dependencies

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