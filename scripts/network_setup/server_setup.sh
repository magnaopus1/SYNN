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

# Function to start the server on the node's IP address
start_server() {
  local node_ip=$1

  echo "Starting server on Node: $NODE_NAME ($NODE_ID) at IP: $node_ip"
  curl -X POST "$node_ip/api/network/server/start" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "node_name": "'"$NODE_NAME"'"
        }'
}

# Main server setup process
setup_server() {
  local node_ip=$1

  # Start the server
  start_server "$node_ip"

  echo "Server setup completed for Node: $NODE_NAME ($NODE_ID) at IP: $node_ip!"
}

# Execute the server setup
setup_server "$NODE_IP"
