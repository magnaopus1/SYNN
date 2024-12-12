#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Bootstrap node IP and connection pool settings
BOOTSTRAP_NODE=$(get_value_from_config "bootstrap_node: ip")
MAX_IDLE_TIME=$(get_value_from_config "bootstrap_node: connection_pool: max_idle_time")
MAX_CONNECTIONS=$(get_value_from_config "bootstrap_node: connection_pool: max_connections")

# New node IP and connection pool settings
NEW_NODE_IP=$(get_value_from_config "new_node_001: ip")
NEW_NODE_MAX_IDLE=$(get_value_from_config "new_node_001: connection_pool: max_idle_time")
NEW_NODE_MAX_CONN=$(get_value_from_config "new_node_001: connection_pool: max_connections")

# Function to set up the connection pool
setup_connection_pool() {
  local node_ip=$1
  local max_idle_time=$2
  local max_connections=$3

  echo "Setting up Connection Pool for Node: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/connection/get_connection" \
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

  echo "Checking active connections for Node: $node_ip"
  curl -X GET "$SERVER_ADDRESS/api/network/connection/active_connections" \
    -H "Content-Type: application/json" \
    -d '{"node_ip": "'"$node_ip"'"}'
}

# Setup connection pool for the bootstrap node
setup_connection_pool "$BOOTSTRAP_NODE" "$MAX_IDLE_TIME" "$MAX_CONNECTIONS"

# Check active connections after setup for the bootstrap node
check_active_connections "$BOOTSTRAP_NODE"

# Setup connection pool for the new node
setup_connection_pool "$NEW_NODE_IP" "$NEW_NODE_MAX_IDLE" "$NEW_NODE_MAX_CONN"

# Check active connections after setup for the new node
check_active_connections "$NEW_NODE_IP"

echo "Connection pool setup completed for Bootstrap Node and New Node!"
