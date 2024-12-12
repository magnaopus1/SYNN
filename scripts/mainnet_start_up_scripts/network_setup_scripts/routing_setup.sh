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

# Function to add a route between nodes
add_route() {
  local from_node=$1
  local to_node=$2

  echo "Adding route from $from_node to $to_node"
  curl -X POST "$SERVER_ADDRESS/api/network/routing/add_route" \
    -H "Content-Type: application/json" \
    -d '{
          "from_node_ip": "'"$from_node"'",
          "to_node_ip": "'"$to_node"'"
        }'
}

# Function to retrieve routes for a node
get_routes() {
  local node_ip=$1

  echo "Getting routes for Node: $node_ip"
  curl -X GET "$SERVER_ADDRESS/api/network/routing/get_route" \
    -H "Content-Type: application/json" \
    -d '{"node_ip": "'"$node_ip"'"}'
}

# Add routes between Bootstrap Node and New Node
add_route "$BOOTSTRAP_NODE" "$NEW_NODE"
add_route "$NEW_NODE" "$BOOTSTRAP_NODE"

# Get routes for both nodes
get_routes "$BOOTSTRAP_NODE"
get_routes "$NEW_NODE"

echo "Routing setup completed!"
