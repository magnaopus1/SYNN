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

# Function to add a route between the current node and the peer
add_route() {
  local from_node=$1
  local to_node=$2

  echo "Adding route from $from_node to $to_node"
  curl -X POST "$from_node/api/network/routing/add_route" \
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
  curl -X GET "$node_ip/api/network/routing/get_route" \
    -H "Content-Type: application/json" \
    -d '{"node_ip": "'"$node_ip"'"}'
}

# Add routes between the current node and the peer
add_route "$NODE_IP" "$PEER_IP"
add_route "$PEER_IP" "$NODE_IP"

# Get routes for the current node and the peer
get_routes "$NODE_IP"
get_routes "$PEER_IP"

echo "Routing setup completed for Node: $NODE_NAME ($NODE_ID) at IP: $NODE_IP!"
