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

# Get node and peer details from the config file
NODE_IP=$(get_value_from_config "node: ip")
PEER_IP=$(get_value_from_config "peer: ip")

# Function to register a node in the network topology
register_node() {
  local node_ip=$1

  echo "Registering Node: $node_ip in the network topology"
  curl -X POST "$node_ip/api/network/topology/register_node" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Function to connect two nodes in the network topology
connect_nodes() {
  local node1_ip=$1
  local node2_ip=$2

  echo "Connecting Node $node1_ip with Node $node2_ip in the topology"
  curl -X POST "$node1_ip/api/network/topology/connect_nodes" \
    -H "Content-Type: application/json" \
    -d '{
          "node1_ip": "'"$node1_ip"'",
          "node2_ip": "'"$node2_ip"'"
        }'
}

# Register the node and its peer in the topology
register_node "$NODE_IP"
if [ -z "$PEER_IP" ]; then
  echo "Peer IP not found. Skipping peer connection."
else
  register_node "$PEER_IP"
  connect_nodes "$NODE_IP" "$PEER_IP"
fi

echo "Network topology setup completed!"
