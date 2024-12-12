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

# Function to add a node to the SDN controller
add_sdn_node() {
  local node_ip=$1

  echo "Adding node to SDN Controller: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/sdn/add_node" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Function to get the status of a node in the SDN controller
get_sdn_node_status() {
  local node_ip=$1

  echo "Getting SDN node status: $node_ip"
  curl -X GET "$SERVER_ADDRESS/api/network/sdn/get_node_status" \
    -H "Content-Type: application/json" \
    -d '{"node_ip": "'"$node_ip"'"}'
}

# Add nodes to SDN Controller
add_sdn_node "$BOOTSTRAP_NODE"
add_sdn_node "$NEW_NODE"

# Get status for both nodes
get_sdn_node_status "$BOOTSTRAP_NODE"
get_sdn_node_status "$NEW_NODE"

echo "SDN Controller setup completed!"
