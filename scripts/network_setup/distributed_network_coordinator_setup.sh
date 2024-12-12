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

# Register new node in the network
register_node() {
  local node_id=$1
  local node_name=$2
  local node_ip=$3

  echo "Registering Node $node_name ($node_id) at $node_ip with the distributed coordinator"
  curl -X POST "$node_ip/api/network/distributed/register_node" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$node_id"'",
          "name": "'"$node_name"'",
          "address": "'"$node_ip"'"
        }'
}

# Establish connection with a peer node
establish_connection() {
  local node_id=$1
  local node_ip=$2

  echo "Establishing connection with Node $node_id"
  curl -X POST "$node_ip/api/network/distributed/establish_connection" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$node_id"'"
        }'
}

# Register and establish connection for the node (using dynamic NODE_IP)
register_node "$NODE_ID" "$NODE_NAME" "$NODE_IP"
establish_connection "$NODE_ID" "$NODE_IP"

echo "Node registration and connection establishment completed for $NODE_TYPE node $NODE_NAME ($NODE_ID) at $NODE_IP!"
