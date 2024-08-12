#!/usr/bin/env bash

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
