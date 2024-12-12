#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Get node IPs
BOOTSTRAP_NODE=$(get_value_from_config "bootstrap_node: ip")
NEW_NODE=$(get_value_from_config "new_node_001: ip")

# Function to discover peers
discover_peers() {
  local node_ip=$1

  echo "Discovering peers for Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/peer_discovery/discover_peers" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Discover peers for both nodes
discover_peers "$BOOTSTRAP_NODE"
discover_peers "$NEW_NODE"

echo "Peer discovery completed!"
