#!/bin/bash

SCRIPTS_PATH="/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/scripts/network_setup"

run_script() {
  local script_name=$1
  echo "Running $script_name..."
  bash "$SCRIPTS_PATH/$script_name"
  if [ $? -ne 0 ]; then
    echo "$script_name failed. Exiting."
    exit 1
  fi
  echo "$script_name completed successfully."
}

run_script "server_setup.sh"

run_script "firewall_setup.sh"

run_script "connection_pool_setup.sh"

run_script "distributed_network_coordinator_setup.sh"

run_script "fault_tolerance_setup.sh"

run_script "setup_peer_discovery.sh"

run_script "geolocation_setup.sh"

run_script "network_setup.sh"

run_script "peer_advertisement_setup.sh"

run_script "p2p_network_setup.sh"

run_script "peer_connections_setup.sh"

run_script "qos_manager_setup.sh"

run_script "routing_setup.sh"

run_script "sdn_controller_setup.sh"

run_script "tls_ssl_setup.sh"

run_script "topology_setup.sh"

echo "Network setup completed successfully!"
