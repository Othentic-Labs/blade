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

- **Geth RPC (Rootchain)**:
  - URL: http://localhost:8545

- **Blockscout Explorer**: 
  - URL: http://localhost:80

## Bridge

The Othentic Blade network implements a cross-chain bridge mechanism to facilitate message passing between Layer 1 (L1) and Layer 2 (L2). This bridge allows for secure and verifiable communication between the two layers, enabling various cross-chain applications and functionalities.

### Bridge Architecture

The bridge consists of several key components:

1. **StateSender Contract**: Deployed on both L1 and L2, this contract initiates the cross-layer message passing.
2. **Receiver Contract**: Deployed on the destination layer (L1 or L2), this contract implements IStateReceiver or IL2StateReceiver accepts and stores the bridged messages.
3. **ExitHelper Contract**: Deployed on L1, this contract facilitates the proof verification and message passing from L2 to L1.

### Bridging Process

#### L1 to L2 Message Passing

1. **Message Initiation**: A transaction is sent to the StateSender contract on L1, including the receiver contract address on L2 and the message data.
2. **Event Emission**: The StateSender contract emits a `StateSynced` event containing a unique state ID, sender address, receiver address, and the message data.
3. **L2 Processing**: The L2 network observes the `StateSynced` event and processes it, delivering the message to the specified receiver contract on L2.
4. **Message Reception**: The receiver contract on L2 stores the received message, emitting a `StateReceived` event with the state ID, sender address, and message data.

#### L2 to L1 Message Passing

1. **Message Initiation**: Similar to L1 to L2, a transaction is sent to the StateSender contract on L2.
2. **Event Emission**: The L2 StateSender contract emits a `StateSynced` event.
3. **Exit Proof Generation**: An exit proof is generated on L2, which includes metadata about the checkpoint block, leaf index, and the exit event, along with a Merkle proof.
4. **Proof Submission**: The exit proof is submitted to the ExitHelper contract on L1.
5. **Verification and Execution**: The ExitHelper contract verifies the proof and executes the message on L1, delivering it to the specified receiver contract.
6. **Message Reception**: The receiver contract on L1 stores the received message and typically emits a `StateReceived` event.

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

## Customization

You can customize the network by modifying the `docker-compose.yml` file and the `blade.sh` script. Key areas for customization include:

- Changing the number of validator nodes
- Adjusting network parameters in the genesis configuration
- Modifying resource allocations for containers
- Adding or removing services

## L2 Contract Deployments

The following are the contract addresses deployed on the L2 network:

### System Contracts
- EpochManagerContract (Proxy): 0x101
- EpochManagerContractV1 (Implementation): 0x1011
- BLSContract (Proxy): 0x102
- BLSContractV1 (Implementation): 0x1021
- MerkleContract (Proxy): 0x103
- MerkleContractV1 (Implementation): 0x1031
- RewardTokenContract (Proxy): 0x104
- RewardTokenContractV1 (Implementation): 0x1041
- DefaultBurnContract (Proxy): 0x106
- StakeManagerContract (Proxy): 0x10022
- StakeManagerContractV1 (Implementation): 0x100221

### Bridge Contracts
- StateReceiverContract (Proxy): 0x1001
- StateReceiverContractV1 (Implementation): 0x10011
- NativeERC20TokenContract (Proxy): 0x1010
- NativeERC20TokenContractV1 (Implementation): 0x10101
- L2StateSenderContract (Proxy): 0x1002
- L2StateSenderContractV1 (Implementation): 0x10021

### Token Contracts
- ChildERC20Contract: 0x1003
- ChildERC20PredicateContract (Proxy): 0x1004
- ChildERC20PredicateContractV1 (Implementation): 0x10041
- ChildERC721Contract: 0x1005
- ChildERC721PredicateContract (Proxy): 0x1006
- ChildERC721PredicateContractV1 (Implementation): 0x10061
- ChildERC1155Contract: 0x1007
- ChildERC1155PredicateContract (Proxy): 0x1008
- ChildERC1155PredicateContractV1 (Implementation): 0x10081
- RootMintableERC20PredicateContract (Proxy): 0x1009
- RootMintableERC20PredicateContractV1 (Implementation): 0x10091
- RootMintableERC721PredicateContract (Proxy): 0x100a
- RootMintableERC721PredicateContractV1 (Implementation): 0x100a1
- RootMintableERC1155PredicateContract (Proxy): 0x100b
- RootMintableERC1155PredicateContractV1 (Implementation): 0x100b1

### Governance Contracts
- ChildGovernorContract (Proxy): 0x100c
- ChildGovernorContractV1 (Implementation): 0x100c1
- ChildTimelockContract (Proxy): 0x100d
- ChildTimelockContractV1 (Implementation): 0x100d1
- NetworkParamsContract (Proxy): 0x100e
- NetworkParamsContractV1 (Implementation): 0x100e1
- ForkParamsContract (Proxy): 0x100f
- ForkParamsContractV1 (Implementation): 0x100f1

### Special Addresses
- SystemCaller: 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE

### Precompiles
- NativeTransferPrecompile: 0x2020
- BLSAggSigsVerificationPrecompile: 0x2030
- ConsolePrecompile: 0x000000000000000000636F6e736F6c652e6c6f67

### Access Control Lists
- AllowListContractsAddr: 0x0200000000000000000000000000000000000000
- BlockListContractsAddr: 0x0300000000000000000000000000000000000000
- AllowListTransactionsAddr: 0x0200000000000000000000000000000000000002
- BlockListTransactionsAddr: 0x0300000000000000000000000000000000000002
- AllowListBridgeAddr: 0x0200000000000000000000000000000000000004
- BlockListBridgeAddr: 0x0300000000000000000000000000000000000004