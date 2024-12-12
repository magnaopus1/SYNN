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
NODE_TYPE=$(get_value_from_config "node: type")
NODE_ID=$(get_value_from_config "node: id")
NODE_NAME=$(get_value_from_config "node: name")

# Get peer details from the config file
PEER_ID=$(get_value_from_config "peer: id")
PEER_IP=$(get_value_from_config "peer: ip")

# Function to start the P2P network
start_p2p_network() {
  local node_ip=$1

  echo "Starting P2P network for Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/p2p/start_network" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "node_name": "'"$NODE_NAME"'",
          "node_type": "'"$NODE_TYPE"'",
          "ip": "'"$node_ip"'"
        }'
}

# Function to connect to a peer in the network
connect_peer() {
  local node_ip=$1
  local peer_ip=$2
  local peer_id=$3

  echo "Connecting to peer $peer_id at IP: $peer_ip from Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/network/connect_peer" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "peer_id": "'"$peer_id"'",
          "peer_ip": "'"$peer_ip"'"
        }'
}

# Function to ping a peer to check connectivity
ping_peer() {
  local node_ip=$1
  local peer_ip=$2
  local peer_id=$3

  echo "Pinging peer $peer_id at IP: $peer_ip from Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/network/ping_peer" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "peer_id": "'"$peer_id"'",
          "peer_ip": "'"$peer_ip"'"
        }'
}

# Main P2P setup process
setup_p2p_network() {
  local node_ip=$1
  local peer_ip=$2
  local peer_id=$3

  # Start the P2P network for the node
  start_p2p_network "$node_ip"

  # Connect to the peer in the network
  connect_peer "$node_ip" "$peer_ip" "$peer_id"

  # Ping the peer to verify connection
  ping_peer "$node_ip" "$peer_ip" "$peer_id"

  echo "P2P network setup completed for Node: $NODE_NAME ($NODE_TYPE, ID: $NODE_ID) at IP: $node_ip!"
}

# Execute the P2P network setup
setup_p2p_network "$NODE_IP" "$PEER_IP" "$PEER_ID"
