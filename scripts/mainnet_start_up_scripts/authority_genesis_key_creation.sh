#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address from config
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Function to create the Genesis Key
create_genesis_key() {
  echo "Creating genesis key for authority node..."
  
  RESPONSE=$(curl -X POST "$SERVER_ADDRESS/api/authority_node/key/genesis_create" \
    -H "Content-Type: application/json")
  
  if [ $? -ne 0 ]; then
    echo "Failed to create genesis key."
    exit 1
  fi

  echo "Genesis key creation response: $RESPONSE"
}

# Main execution
create_genesis_key

echo "Genesis Key creation completed."
