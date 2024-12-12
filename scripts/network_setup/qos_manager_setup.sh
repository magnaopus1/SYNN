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
PEER_IP=$(get_value_from_config "peer: ip")

# Function to add the node to the QoS manager
add_peer_qos() {
  local node_ip=$1

  echo "Adding Node $node_ip to QoS Manager"
  curl -X POST "$node_ip/api/network/qos/add_peer" \
    -H "Content-Type: application/json" \
    -d '{
          "node_ip": "'"$node_ip"'"
        }'
}

# Function to send a test packet to verify QoS
send_qos_packet() {
  local from_ip=$1
  local to_ip=$2

  echo "Sending QoS packet from $from_ip to $to_ip"
  curl -X POST "$from_ip/api/network/qos/send_packet" \
    -H "Content-Type: application/json" \
    -d '{
          "from_ip": "'"$from_ip"'",
          "to_ip": "'"$to_ip"'",
          "packet_size": 1024,          # Sample packet size (in bytes)
          "priority": "high"            # Priority level for the packet
        }'
}

# Add the node and peer to QoS manager
add_peer_qos "$NODE_IP"
add_peer_qos "$PEER_IP"

# Test QoS by sending a packet between the node and peer
send_qos_packet "$NODE_IP" "$PEER_IP"

echo "QoS Manager setup completed for Node: $NODE_NAME ($NODE_ID) at IP: $NODE_IP!"
