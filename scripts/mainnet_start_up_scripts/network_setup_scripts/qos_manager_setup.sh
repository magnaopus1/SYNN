#!/bin/bash

CONFIG_FILE="/path/to/network_config.yaml"

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Get server address
SERVER_ADDRESS=$(get_value_from_config "server: address")

# Get node IPs
BOOTSTRAP_NODE=$(get_value_from_config "bootstrap_node: ip")
NEW_NODE=$(get_value_from_config "new_node_001: ip")

# Function to add a peer to the QoS management system
add_peer_qos() {
  local node_ip=$1

  echo "Adding Peer to QoS Manager: $node_ip"
  curl -X POST "$SERVER_ADDRESS/api/network/qos/add_peer" \
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
  curl -X POST "$SERVER_ADDRESS/api/network/qos/send_packet" \
    -H "Content-Type: application/json" \
    -d '{
          "from_ip": "'"$from_ip"'",
          "to_ip": "'"$to_ip"'",
          "packet_size": 1024,          # Sample packet size (in bytes)
          "priority": "high"            # Priority level for the packet
        }'
}

# Add both nodes to QoS manager
add_peer_qos "$BOOTSTRAP_NODE"
add_peer_qos "$NEW_NODE"

# Test QoS by sending a packet between the nodes
send_qos_packet "$BOOTSTRAP_NODE" "$NEW_NODE"

echo "QoS Manager setup completed!"
