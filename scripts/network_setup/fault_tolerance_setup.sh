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

# Function to handle node failure
node_failure() {
  local node_ip=$1

  echo "Handling Node Failure for Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/node_failure" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "node_name": "'"$NODE_NAME"'",
          "node_type": "'"$NODE_TYPE"'",
          "ip": "'"$node_ip"'"
        }'
}

# Function to handle node recovery
node_recovery() {
  local node_ip=$1

  echo "Handling Node Recovery for Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/node_recovery" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "node_name": "'"$NODE_NAME"'",
          "node_type": "'"$NODE_TYPE"'",
          "ip": "'"$node_ip"'"
        }'
}

# Function to check if the quorum is alive
is_quorum_alive() {
  local node_ip=$1

  echo "Checking if Quorum is Alive for Node: $NODE_NAME ($NODE_ID)"
  curl -X GET "$node_ip/api/network/is_quorum_alive" \
    -H "Content-Type: application/json"
}

# Function to sync node with the network
sync_node() {
  local node_ip=$1

  echo "Syncing Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/sync_node" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "node_name": "'"$NODE_NAME"'",
          "node_type": "'"$NODE_TYPE"'",
          "ip": "'"$node_ip"'"
        }'
}

# Main fault tolerance handler
handle_fault_tolerance() {
  local node_ip=$1

  echo "Starting Fault Tolerance Setup for Node: $NODE_NAME ($NODE_ID)"

  # Handle node failure
  node_failure "$node_ip"

  # Check quorum status
  is_quorum_alive "$node_ip"

  # If quorum is alive, attempt recovery and sync
  node_recovery "$node_ip"
  sync_node "$node_ip"

  echo "Fault Tolerance Setup Completed for Node: $NODE_NAME ($NODE_TYPE, ID: $NODE_ID)"
}

# Execute the fault tolerance process
handle_fault_tolerance "$NODE_IP"
