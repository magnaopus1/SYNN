#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address from config
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Step 1: Create a new wallet
create_new_wallet() {
  echo "Creating new wallet..."
  WALLET_ID=$(curl -X POST "$SERVER_ADDRESS/api/wallet/new_wallet" \
    -H "Content-Type: application/json" | jq -r '.wallet_id')

  if [ -z "$WALLET_ID" ]; then
    echo "Failed to create wallet"
    exit 1
  fi
  echo "Wallet created with ID: $WALLET_ID"
}

# Step 2: Assign name 'genesis_wallet' to the new wallet
assign_wallet_name() {
  echo "Assigning name 'genesis_wallet' to wallet..."
  curl -X POST "$SERVER_ADDRESS/api/wallet/assign_name_to_wallet" \
    -H "Content-Type: application/json" \
    -d '{"wallet_id": "'"$WALLET_ID"'", "wallet_name": "genesis_wallet"}'

  echo "Wallet name 'genesis_wallet' assigned."
}

# Step 3: Create backup for the wallet
create_wallet_backup() {
  echo "Creating backup for wallet..."
  curl -X POST "$SERVER_ADDRESS/api/wallet/create_backup" \
    -H "Content-Type: application/json" \
    -d '{"wallet_id": "'"$WALLET_ID"'"}'

  echo "Backup created for wallet."
}

# Step 4: Generate mnemonic for recovery
generate_mnemonic() {
  echo "Generating mnemonic for wallet recovery..."
  MNEMONIC=$(curl -X POST "$SERVER_ADDRESS/api/wallet/generate_mnemonic" \
    -H "Content-Type: application/json" \
    -d '{"wallet_id": "'"$WALLET_ID"'"}' | jq -r '.mnemonic')

  if [ -z "$MNEMONIC" ]; then
    echo "Failed to generate mnemonic"
    exit 1
  fi
  echo "Mnemonic generated: $MNEMONIC"
}

# Step 5: Display balance for the genesis wallet
display_wallet_balance() {
  echo "Displaying wallet balance..."
  curl -X GET "$SERVER_ADDRESS/api/wallet/display_balances" \
    -H "Content-Type: application/json" \
    -d '{"wallet_id": "'"$WALLET_ID"'"}'

  echo "Wallet balance displayed."
}

# Main execution flow
create_new_wallet
assign_wallet_name
create_wallet_backup
generate_mnemonic
display_wallet_balance

echo "Genesis Wallet setup complete."
