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

# Function to initiate TLS handshake
start_tls_handshake() {
  local node_ip=$1

  echo "Starting TLS handshake for Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/tls/start_tls_handshake" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Function to initiate SSL handshake
start_ssl_handshake() {
  local node_ip=$1

  echo "Starting SSL handshake for Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/ssl/start_handshake" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Start TLS and SSL handshakes for both nodes
start_tls_handshake "$BOOTSTRAP_NODE"
start_ssl_handshake "$BOOTSTRAP_NODE"

start_tls_handshake "$NEW_NODE"
start_ssl_handshake "$NEW_NODE"

echo "TLS and SSL handshakes have been initiated for all nodes!"
