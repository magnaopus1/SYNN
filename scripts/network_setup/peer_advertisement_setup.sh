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
NODE_ID=$(get_value_from_config "node: id")
NODE_NAME=$(get_value_from_config "node: name")
NODE_TYPE=$(get_value_from_config "node: type")

# Function to advertise the peer node on the network
advertise_peer() {
  local node_ip=$1

  echo "Advertising Node: $NODE_NAME ($NODE_ID) at IP: $node_ip"
  curl -X POST "$node_ip/api/network/peer_advertiser/advertise_peer" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "node_name": "'"$NODE_NAME"'",
          "node_type": "'"$NODE_TYPE"'",
          "node_ip": "'"$node_ip"'"
        }'
}

# Advertise the node as a peer
advertise_peer "$NODE_IP"

echo "Peer advertisement setup completed for Node: $NODE_NAME ($NODE_TYPE, ID: $NODE_ID) at IP: $NODE_IP!"
