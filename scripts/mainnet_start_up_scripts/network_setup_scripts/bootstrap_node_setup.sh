#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Bootstrap node configuration (using only one bootstrap node now)
BOOTSTRAP_NODE=$(get_value_from_config "bootstrap_node: ip")

# Function to start the bootstrap node using API call
start_bootstrap_node() {
  local node_ip=$1
  echo "Starting Bootstrap Node at $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/bootstrap/start_node" \
    -H "Content-Type: application/json" \
    -d '{"ip": "'"$node_ip"'"}'
}

# Start Bootstrap Node
start_bootstrap_node "$BOOTSTRAP_NODE"

echo "Bootstrap Node has been started!"
