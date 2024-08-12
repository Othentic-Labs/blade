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

# Function to check if all services are healthy
check_services_health() {
    local timeout=300  # 5 minutes timeout
    local start_time=$(date +%s)
    
    while true; do
        if docker-compose ps | grep -q "unhealthy"; then
            local current_time=$(date +%s)
            if [ $((current_time - start_time)) -ge $timeout ]; then
                echo "Timeout: Some services are still unhealthy after 5 minutes."
                return 1
            fi
            sleep 5
        else
            echo "All services are healthy!"
            return 0
        fi
    done
}

# Function to show available services message
show_services_message() {
    echo "
The following services are now available:

* Blade JSON-RPC:
  http://localhost:10002
* Blockscout Explorer:
  http://localhost:80
* Geth RPC (Rootchain):
  http://localhost:8545
"
}

# Function to display ASCII art
show_ascii_art() {
    echo "
  ____  _   _                 _   _       ____  _           _      
 / __ \| | | |               | | (_)     |  _ \| |         | |     
| |  | | |_| |__   ___ _ __ | |_ _  ___  | |_) | | __ _  __| | ___ 
| |  | | __| '_ \ / _ \ '_ \| __| |/ __| |  _ <| |/ _` |/ _` |/ _ \\
| |__| | |_| | | |  __/ | | | |_| | (__  | |_) | | (_| | (_| |  __/
 \____/ \__|_| |_|\___|_| |_|\__|_|\___| |____/|_|\__,_|\__,_|\___|
"
}

# Check dependencies before proceeding
check_dependencies

case "$1" in
    --destroy)
        destroy_environment
        ;;
    --local)
        show_ascii_art
        echo "Running local environment..."
        cd docker/local
        docker-compose up -d
        if check_services_health; then
            show_services_message
        fi
        ;;
    --holesky)
        show_ascii_art
        echo "Running Holesky environment..."
        cd docker/local
        docker-compose -f docker-compose_holesky.yml up -d
        if check_services_health; then
            show_services_message
        fi
        ;;
    --help|*)
        show_help
        ;;
esac