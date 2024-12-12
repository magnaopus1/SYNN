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

# Function to add a node to the SDN controller
add_sdn_node() {
  local node_ip=$1

  echo "Adding Node $node_ip to SDN Controller"
  curl -X POST "$node_ip/api/network/sdn/add_node" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Function to get the status of a node in the SDN controller
get_sdn_node_status() {
  local node_ip=$1

  echo "Getting SDN Node Status for $node_ip"
  curl -X GET "$node_ip/api/network/sdn/get_node_status" \
    -H "Content-Type: application/json" \
    -d '{"node_ip": "'"$node_ip"'"}'
}

# Add the node and its peer to the SDN controller
add_sdn_node "$NODE_IP"
add_sdn_node "$PEER_IP"

# Get status for both the node and the peer
get_sdn_node_status "$NODE_IP"
get_sdn_node_status "$PEER_IP"

echo "SDN Controller setup completed for Node: $NODE_NAME ($NODE_ID) at IP: $NODE_IP!"
