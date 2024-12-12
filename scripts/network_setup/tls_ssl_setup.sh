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
PEER_IP=$(get_value_from_config "peer: ip")

# Function to initiate TLS handshake
start_tls_handshake() {
  local node_ip=$1

  echo "Starting TLS handshake for Node: $node_ip"
  curl -X POST "$node_ip/api/network/tls/start_tls_handshake" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Function to initiate SSL handshake
start_ssl_handshake() {
  local node_ip=$1

  echo "Starting SSL handshake for Node: $node_ip"
  curl -X POST "$node_ip/api/network/ssl/start_handshake" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Start TLS and SSL handshakes for the node
start_tls_handshake "$NODE_IP"
start_ssl_handshake "$NODE_IP"

# Start TLS and SSL handshakes for the peer
if [ -z "$PEER_IP" ]; then
  echo "Peer IP not found. Skipping peer handshake."
else
  start_tls_handshake "$PEER_IP"
  start_ssl_handshake "$PEER_IP"
fi

echo "TLS and SSL handshakes have been initiated for Node and Peer!"
