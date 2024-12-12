#!/bin/bash

echo "Starting Mainnet Network Setup..."

# Start Bootstrap Node
echo "Step 1: Starting Bootstrap Node..."
./bootstrap_node_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to start Bootstrap Node. Exiting setup."
  exit 1
fi

# Set up Connection Pool
echo "Step 2: Setting up Connection Pool..."
./connection_pool_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up Connection Pool. Exiting setup."
  exit 1
fi

# Register Nodes and Set up Distributed Network Coordination
echo "Step 3: Registering Nodes and Setting up Distributed Network Coordination..."
./distributed_network_coordinator_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up Distributed Network Coordination. Exiting setup."
  exit 1
fi

# Set up Fault Tolerance
echo "Step 4: Setting up Fault Tolerance Mechanisms..."
./fault_tolerance_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up Fault Tolerance. Exiting setup."
  exit 1
fi

# Set up Firewall Rules
echo "Step 5: Configuring Firewall Rules..."
./firewall_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up Firewall. Exiting setup."
  exit 1
fi

# Set up QoS Manager
echo "Step 6: Setting up QoS Manager..."
./qos_manager_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up QoS Manager. Exiting setup."
  exit 1
fi

# Set up Routing
echo "Step 7: Setting up Routing..."
./routing_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up Routing. Exiting setup."
  exit 1
fi

# Set up SDN Controller
echo "Step 8: Configuring SDN Controller..."
./sdn_controller_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up SDN Controller. Exiting setup."
  exit 1
fi

# Start Blockchain Server
echo "Step 9: Starting Blockchain Server..."
./server_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to start Blockchain Server. Exiting setup."
  exit 1
fi

# Set up Peer Advertisement
echo "Step 10: Setting up Peer Advertisement..."
./peer_advertisement_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up Peer Advertisement. Exiting setup."
  exit 1
fi

# Discover Peers
echo "Step 11: Discovering Peers..."
./setup_peer_discovery.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to discover peers. Exiting setup."
  exit 1
fi

# Set up TLS and SSL Handshake
echo "Step 12: Setting up TLS and SSL Handshake..."
./tls_ssl_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up TLS/SSL Handshake. Exiting setup."
  exit 1
fi

# Set up Network Topology
echo "Step 13: Setting up Network Topology..."
./topology_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to set up Network Topology. Exiting setup."
  exit 1
fi

# Initialize P2P Network
echo "Step 14: Initializing P2P Network..."
./p2p_network_setup.sh
if [ $? -ne 0 ]; then
  echo "Error: Failed to initialize P2P Network. Exiting setup."
  exit 1
fi

echo "Mainnet Network Setup Complete!"
