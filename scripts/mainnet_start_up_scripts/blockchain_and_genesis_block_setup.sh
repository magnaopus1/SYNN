#!/bin/bash

# Path to the configuration file
CONFIG_PATH="../config/mainnet_config.yaml"

# API endpoint for blockchain initialization
BLOCKCHAIN_INIT_URL="http://localhost:8080/api/blockchain/init"

# Initialize the blockchain and genesis block
echo "Initializing Blockchain and Genesis Block..."

# Make the request to the blockchain init API
response=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BLOCKCHAIN_INIT_URL)

# Check if the request was successful
if [ "$response" -eq 200 ]; then
    echo "Blockchain and Genesis Block successfully initialized."
else
    echo "Failed to initialize Blockchain. HTTP status code: $response"
    exit 1
fi

# Check if the blockchain is up and running
echo "Verifying blockchain initialization..."
response=$(curl -s -o /dev/null -w "%{http_code}" -X GET "http://localhost:8080/api/blockchain/status")

if [ "$response" -eq 200 ]; then
    echo "Blockchain is up and running."
else
    echo "Blockchain initialization failed. Status check returned HTTP code: $response"
    exit 1
fi

# Script finished
echo "Blockchain and Genesis Block setup completed."
