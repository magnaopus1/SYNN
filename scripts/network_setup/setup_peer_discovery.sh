#!/bin/bash

# Prompt user for config file path if not provided as an argument
if [ -z "$1" ]; then
  read -p "Please enter the configuration file path (or press Enter to use the default): " CONFIG_FILE
  CONFIG_FILE=${CONFIG_FILE:-"/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/configs/standard_configs/node_config.yaml"}
else
  CONFIG_FILE=$1
fi

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get node details from the config file
NODE_IP=$(get_value_from_config "node: ip")
NODE_ID=$(get_value_from_config "node: id")
NODE_NAME=$(get_value_from_config "node: name")

# Get peer ID and peer IP from the config file, if available
PEER_ID=$(get_value_from_config "peer: id")
PEER_IP=$(get_value_from_config "peer: ip")

# Function to discover peers dynamically
discover_peers() {
  local node_ip=$1

  echo "Discovering peers for Node: $NODE_NAME ($NODE_ID) at IP: $node_ip"
  curl -X POST "$node_ip/api/network/peer_discovery/discover_peers" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Check if peer ID or peer IP is provided
if [ -z "$PEER_ID" ] || [ -z "$PEER_IP" ]; then
  echo "Peer ID or Peer IP not provided in the config. Discovering peers dynamically..."
  discover_peers "$NODE_IP"
else
  echo "Peer ID: $PEER_ID and Peer IP: $PEER_IP found in the config."
fi

echo "Peer discovery/setup completed for Node: $NODE_NAME ($NODE_ID)!"
