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
PEER_IP=$(get_value_from_config "peer: ip")

# Function to connect to a peer
connect_peer() {
  local node_ip=$1
  local peer_ip=$2

  echo "Connecting to Peer: $peer_ip from Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/peer_connection/connect" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "peer_ip": "'"$peer_ip"'"
        }'
}

# Function to disconnect from a peer
disconnect_peer() {
  local node_ip=$1
  local peer_ip=$2

  echo "Disconnecting from Peer: $peer_ip for Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/peer_connection/disconnect" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "peer_ip": "'"$peer_ip"'"
        }'
}

# Main peer connection setup process
setup_peer_connections() {
  local node_ip=$1
  local peer_ip=$2

  # Connect to the peer
  connect_peer "$node_ip" "$peer_ip"

  echo "Peer connection setup completed for Node: $NODE_NAME ($NODE_ID) at IP: $node_ip!"
}

# Execute the peer connection setup
setup_peer_connections "$NODE_IP" "$PEER_IP"
