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

# Function to register a node in the network topology
register_node() {
  local node_ip=$1

  echo "Registering Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/topology/register_node" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Function to connect two nodes in the network topology
connect_nodes() {
  local node1_ip=$1
  local node2_ip=$2

  echo "Connecting Node $node1_ip with Node $node2_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/topology/connect_nodes" \
    -H "Content-Type: application/json" \
    -d '{
          "node1_ip": "'"$node1_ip"'",
          "node2_ip": "'"$node2_ip"'"
        }'
}

# Register both nodes in the topology
register_node "$BOOTSTRAP_NODE"
register_node "$NEW_NODE"

# Connect the nodes
connect_nodes "$BOOTSTRAP_NODE" "$NEW_NODE"

echo "Network topology setup completed!"
