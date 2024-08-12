# Othentic Blade

Othentic Blade is a fork of Polygon Edge, providing a customizable blockchain network solution. This README will guide you through setting up and running the Othentic Blade network.

## Quick Start

In the root of the repository, you'll find a `run.sh` script that helps you manage the Othentic Blade environment.

### Usage

```bash
./run.sh [OPTIONS]
```

Options:
- `--destroy`: Destroy the environment
- `--local`: Run Blade network with local configuration
- `--holesky`: Run Blade network with Holesky L1
- `--help`: Display the help message

## Script Flow

The `run.sh` script, along with the Docker Compose configuration, sets up the following components:

1. **Local Geth or Holesky Network (L1)**
   - Runs as the rootchain for the Blade network
   - Exposes RPC on port 8545

2. **Blade Network**
   - Consists of 4 validator nodes
   - Each node exposes the following ports:
     - gRPC: 9632
     - JSON-RPC: 8545
     - Prometheus: 5001
     - LibP2P: 1478

3. **Blockscout Block Explorer**
   - Provides a web interface to explore the Blade network
   - Includes supporting services:
     - PostgreSQL database
     - Redis
     - Blockscout backend
     - Blockscout frontend

### Detailed Component Breakdown

#### Rootchain
- Uses the `ghcr.io/0xpolygon/go-ethereum-console:latest` image
- Runs in dev mode with a 2-second block time
- Exposes HTTP and WebSocket interfaces

#### Blade Network Initialization
- Uses a custom Dockerfile to build the Blade binary
- Initializes the network configuration and generates the genesis file

#### Validator Nodes (1-4)
- Run the Blade server with specific configurations
- Each node has its own data directory and sealing enabled

#### Blockscout Services
- **Database**: PostgreSQL for storing blockchain data
- **Redis**: Used for caching and temporary storage
- **Backend**: Processes and indexes blockchain data
- **Frontend**: Serves the web interface for blockchain exploration
- **Stats**: Collects and serves network statistics
- **Visualizer**: Provides network visualization capabilities
- **Proxy**: NGINX server to route traffic to appropriate services

## Running the Network

1. To start the local Blade network:
   ```bash
   ./run.sh --local
   ```

2. To start the Blade network with Holesky L1:
   ```bash
   ./run.sh --holesky
   ```

3. To destroy the environment and clean up:
   ```bash
   ./run.sh --destroy
   ```

## Accessing Services

- **Blade JSON-RPC**:
  - Node 1: http://localhost:10002
  - Node 2: http://localhost:20002
  - Node 3: http://localhost:30002
  - Node 4: http://localhost:40002

- **Blockscout Explorer**: 
  - URL: http://localhost:80

- **Geth RPC (Rootchain)**:
  - URL: http://localhost:8545

## Development

The Dockerfile in the `docker/local` directory is used to build the Blade binary. It uses a multi-stage build process:

1. **Builder Stage**:
   - Uses `golang:1.21-alpine` as the base image
   - Builds the Blade binary using the project's Makefile

2. **Runner Stage**:
   - Uses `alpine:latest` as the base image
   - Copies the built binary and necessary scripts
   - Sets up the environment for running the Blade network

## Customization

You can customize the network by modifying the `docker-compose.yml` file and the `blade.sh` script. Key areas for customization include:

- Changing the number of validator nodes
- Adjusting network parameters in the genesis configuration
- Modifying resource allocations for containers
- Adding or removing services

## Troubleshooting

If you encounter issues:

1. Check the logs of individual services using `docker-compose logs [service_name]`
2. Ensure all required ports are free on your host machine
3. Verify that the rootchain (local Geth or Holesky) is running and accessible

For more detailed information, please refer to the original Polygon Edge documentation and adapt it to the specifics of Othentic Blade.