#!/bin/sh

set -e

echo "Waiting for rootchain to be ready..."
sleep 30

BLADE_BIN=./blade
CHAIN_CUSTOM_OPTIONS=$(tr "\n" " " << EOL
--block-gas-limit 10000000
--epoch-size 10
--chain-id 51001
--name Blade
--premine 0x0000000000000000000000000000000000000000:0xD3C21BCECCEDA1000000
EOL
)

# Deploy the ERC20 token using Node.js script
deployERC20() {
  echo "Deploying ERC20 token on the rootchain..."
  node deploy_erc20.js \
    "http://rootchain:8545" \
    "71394bcfed7228c0a33a9e65b42ba8ce8a697ffdf15dc862aa7aece27819e938" \
    "DanielToken" \
    "DAN" \
    18 \
    "100000000000000000000000000" \
    "$addresses"

  ERC20_ADDRESS=$(cat /data/erc20_address.txt)
  echo "ERC20 token deployed at address: $ERC20_ADDRESS"
}

# createGenesisConfig creates genesis configuration
createGenesisConfig() {
  local consensus_type="$1"
  local secrets="$2"
  shift 2
  echo "Generating $consensus_type Genesis file..."

  "$BLADE_BIN" genesis $CHAIN_CUSTOM_OPTIONS \
    --dir /data/genesis.json \
    --validators-path /data \
    --validators-prefix data- \
    --consensus $consensus_type \
    --bootnode "/dns4/node-1/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[0] | .node_id')" \
    --bootnode "/dns4/node-2/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[1] | .node_id')" \
    --bootnode "/dns4/node-3/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[2] | .node_id')" \
    --bootnode "/dns4/node-4/tcp/1478/p2p/$(echo "$secrets" | jq -r '.[3] | .node_id')" \
    "$@"
}

case "$1" in
   "init")
      case "$2" in
          "polybft")
              echo "Generating PolyBFT secrets..."
              secrets=$("$BLADE_BIN" secrets init --insecure --num 4 --data-dir /data/data- --json)
              echo "Secrets have been successfully generated"

              rm -f /data/genesis.json

              proxyContractsAdmin=0x80cd9D056bc38ECA50cF74A7b9F4d0FB897152a2

              addresses="$(echo "$secrets" | jq -r '.[0] | .address'),$(echo "$secrets" | jq -r '.[1] | .address'),$(echo "$secrets" | jq -r '.[2] | .address'),$(echo "$secrets" | jq -r '.[3] | .address')"
              
              echo "Addresses: $addresses"

              "$BLADE_BIN" bridge fund \
                --json-rpc http://rootchain:8545 \
                --addresses ${proxyContractsAdmin} \
                --amounts 1000000000000000000000000
              
              deployERC20

              createGenesisConfig "$2" "$secrets" \
                --reward-wallet 0xDEADBEEF:1000000 \
                --native-token-config "Blade:BLD:18:false" \
                --blade-admin $(echo "$secrets" | jq -r '.[0] | .address') \
                --proxy-contracts-admin ${proxyContractsAdmin}

              echo "ERC20 token address: $ERC20_ADDRESS"

              "$BLADE_BIN" bridge deploy \
                --json-rpc http://rootchain:8545 \
                --genesis /data/genesis.json \
                --proxy-contracts-admin ${proxyContractsAdmin} \
                --erc20-token ${ERC20_ADDRESS} \
                --test


              "$BLADE_BIN" bridge fund \
                --json-rpc http://rootchain:8545 \
                --addresses ${addresses} \
                --amounts 1000000000000000000000000,1000000000000000000000000,1000000000000000000000000,1000000000000000000000000

              ;;
      esac
      ;;
  "start-node-1")
    relayer_flag=""
    # Start relayer only when run in polybft
    if [ "$2" == "polybft" ]; then
      echo "Starting relayer..."
      relayer_flag="--relayer"
    fi

    "$BLADE_BIN" server \
      --data-dir /data/data-1 \
      --chain /data/genesis.json \
      --grpc-address 0.0.0.0:9632 \
      --libp2p 0.0.0.0:1478 \
      --jsonrpc 0.0.0.0:8545 \
      --prometheus 0.0.0.0:5001 \
      $relayer_flag
   ;;
   *)
      echo "Executing blade..."
      exec "$BLADE_BIN" "$@"
      ;;
esac
