#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Initialize Charity Pool
initialize_charity_pool() {
  echo "Initializing Charity Pool..."
  curl -X POST "$SERVER_ADDRESS/api/charity/init_pool" \
    -H "Content-Type: application/json"
  echo "Charity Pool initialized."
}

# Initialize External Charity Pool
initialize_external_charity_pool() {
  echo "Initializing External Charity Pool..."
  curl -X POST "$SERVER_ADDRESS/api/charity/init_external_pool" \
    -H "Content-Type: application/json"
  echo "External Charity Pool initialized."
}

# Initialize Internal Charity Pool
initialize_internal_charity_pool() {
  echo "Initializing Internal Charity Pool..."
  curl -X POST "$SERVER_ADDRESS/api/charity/init_internal_pool" \
    -H "Content-Type: application/json"
  echo "Internal Charity Pool initialized."
}

# Call the initialization functions
initialize_charity_pool
initialize_external_charity_pool
initialize_internal_charity_pool

echo "All charity pools have been initialized successfully."
