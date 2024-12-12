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

# Function to advertise peer on the network
advertise_peer() {
  local node_ip=$1

  echo "Advertising Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/peer_advertiser/advertise_peer" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Advertise both nodes on the network
advertise_peer "$BOOTSTRAP_NODE"
advertise_peer "$NEW_NODE"

echo "Peer advertisement setup completed!"
