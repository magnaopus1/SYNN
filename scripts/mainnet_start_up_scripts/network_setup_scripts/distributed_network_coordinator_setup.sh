#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Bootstrap node IP
BOOTSTRAP_NODE=$(get_value_from_config "bootstrap_node: ip")

# Register new nodes in the network
register_node() {
  local node_id=$1
  local node_ip=$2

  echo "Registering Node $node_id at $node_ip with the distributed coordinator"
  curl -X POST "$SERVER_ADDRESS/api/network/distributed/register_node" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$node_id"'",
          "address": "'"$node_ip"'"
        }'
}

# Establish connection with a peer node
establish_connection() {
  local node_id=$1

  echo "Establishing connection with Node $node_id"
  curl -X POST "$SERVER_ADDRESS/api/network/distributed/establish_connection" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$node_id"'"
        }'
}

# Register and establish connection for the new node
register_node "new_node_001" "http://192.168.1.105:8080"
establish_connection "new_node_001"

echo "Node registration and connection establishment completed!"
