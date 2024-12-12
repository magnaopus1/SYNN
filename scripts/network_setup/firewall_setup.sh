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

# Get firewall settings
FIREWALL_ENABLED=$(get_value_from_config "firewall: enabled")

# Get whitelisted and blacklisted IPs
WHITELISTED_IPS=$(get_value_from_config "firewall: whitelisted_ips")
BLACKLISTED_IPS=$(get_value_from_config "firewall: blacklisted_ips")

# Function to allow an IP address in the firewall
allow_ip() {
  local node_ip=$1
  local ip_to_allow=$2

  echo "Allowing IP $ip_to_allow for Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/firewall/allow_ip" \
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

  echo "Blocking IP $ip_to_block for Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/firewall/block_ip" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'",
          "ip_to_block": "'"$ip_to_block"'"
        }'
}

# Process whitelisted and blacklisted IPs if firewall is enabled
if [ "$FIREWALL_ENABLED" == "true" ]; then
  echo "Firewall is enabled for Node: $NODE_NAME ($NODE_ID)"

  # Allow whitelisted IPs
  for IP in $(echo $WHITELISTED_IPS | tr "," "\n"); do
    allow_ip "$NODE_IP" "$IP"
  done

  # Block blacklisted IPs
  for IP in $(echo $BLACKLISTED_IPS | tr "," "\n"); do
    block_ip "$NODE_IP" "$IP"
  done
else
  echo "Firewall is disabled for Node: $NODE_NAME ($NODE_ID)"
fi

# If the IP is neither in whitelist nor blacklist, it's allowed by default
if [ -z "$WHITELISTED_IPS" ] && [ -z "$BLACKLISTED_IPS" ]; then
  echo "No whitelist or blacklist found. All IPs allowed by default for Node: $NODE_NAME ($NODE_ID)."
fi

echo "Firewall setup completed for Node: $NODE_NAME ($NODE_TYPE, ID: $NODE_ID) at IP: $NODE_IP!"
