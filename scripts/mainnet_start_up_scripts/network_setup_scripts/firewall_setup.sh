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

# Function to allow an IP address in the firewall
allow_ip() {
  local node_ip=$1
  local ip_to_allow=$2

  echo "Allowing IP $ip_to_allow for Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/firewall/allow_ip" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'",
          "ip_to_allow": "'"$ip_to_allow"'"
        }'
}

# Function to block an IP address in the firewall
block_ip() {
  local node_ip=$1
  local ip_to_block=$2

  echo "Blocking IP $ip_to_block for Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/firewall/block_ip" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'",
          "ip_to_block": "'"$ip_to_block"'"
        }'
}

# Allow necessary IPs for both nodes (example IPs to allow)
allow_ip "$BOOTSTRAP_NODE" "192.168.1.200"
allow_ip "$NEW_NODE" "192.168.1.201"

# Block any untrusted IPs (example IPs to block)
block_ip "$BOOTSTRAP_NODE" "192.168.1.50"
block_ip "$NEW_NODE" "192.168.1.51"

echo "Firewall rules have been set!"
