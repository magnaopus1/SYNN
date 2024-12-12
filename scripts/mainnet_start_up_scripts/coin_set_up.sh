#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address from the config file
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Coin initialization API call
echo "Initializing the coin supply and setup..."

# Since the coin supply and details are hardcoded in the function, we just need to trigger the initialization.
curl -X POST "$SERVER_ADDRESS/api/coin/init_supply" \
  -H "Content-Type: application/json" \
  -d '{}'

# Check for success response
if [ $? -eq 0 ]; then
    echo "Coin supply initialization completed successfully!"
else
    echo "Error occurred during coin initialization."
fi

echo "Coin setup process completed."
