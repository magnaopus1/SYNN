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
NODE_TYPE=$(get_value_from_config "node: type")
NODE_ID=$(get_value_from_config "node: id")
NODE_NAME=$(get_value_from_config "node: name")

# Get connection pool settings from the config file
MAX_IDLE_TIME=$(get_value_from_config "node: connection_pool: max_idle_time")
MAX_CONNECTIONS=$(get_value_from_config "node: connection_pool: max_connections")

# Function to set up the connection pool
setup_connection_pool() {
  local node_ip=$1
  local max_idle_time=$2
  local max_connections=$3

  echo "Setting up Connection Pool for Node $NODE_NAME ($NODE_ID) at IP: $node_ip"
  curl -X POST "$node_ip/api/network/connection/get_connection" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'",
          "max_idle_time": '"$max_idle_time"',
          "max_connections": '"$max_connections"'
        }'
}

# Function to check active connections for the node
check_active_connections() {
  local node_ip=$1

  echo "Checking active connections for Node $NODE_NAME ($NODE_ID) at IP: $node_ip"
  curl -X GET "$node_ip/api/network/connection/active_connections" \
    -H "Content-Type: application/json" \
    -d '{"node_ip": "'"$node_ip"'"}'
}

# Setup connection pool for the node
setup_connection_pool "$NODE_IP" "$MAX_IDLE_TIME" "$MAX_CONNECTIONS"

# Check active connections after setup
check_active_connections "$NODE_IP"

echo "Connection pool setup completed for Node: $NODE_NAME ($NODE_TYPE, ID: $NODE_ID) at IP: $NODE_IP!"
