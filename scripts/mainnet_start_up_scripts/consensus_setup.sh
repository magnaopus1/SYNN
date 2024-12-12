#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Function to start the first iteration of the consensus process
start_consensus() {
  echo "Starting the first iteration of the consensus process..."
  curl -X POST "$SERVER_ADDRESS/api/consensus/process_transactions" \
    -H "Content-Type: application/json"
  echo "Consensus process started and transactions are being processed."
}

# Call the function to start the consensus process
start_consensus

echo "First consensus iteration complete. Further consensus will be handled by the consensus loop."
