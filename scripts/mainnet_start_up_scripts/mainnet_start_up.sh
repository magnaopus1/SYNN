#!/bin/bash

# Function to display a progress bar
progress_bar() {
    local duration=$1
    already_done() { for ((done=0; done<$elapsed; done++)); do printf "#"; done }
    remaining() { for ((remain=$elapsed; remain<$duration; remain++)); do printf "."; done }
    percentage() { printf "| %s%%" $(( ($elapsed*100)/($duration*1) )); }

    for ((elapsed=1; elapsed<=duration; elapsed++)); do
        already_done; remaining; percentage
        sleep 1
        printf "\r"
    done
    printf "\n"
}

# Function to run a script with logging and progress bar
run_script() {
    script_path=$1
    log_file=$2
    echo "Running $script_path ..."
    bash $script_path >> $log_file 2>&1
    if [ $? -ne 0 ]; then
        echo "Error while running $script_path. Check the log: $log_file"
        exit 1
    fi
    echo "$script_path completed."
}

# Mainnet Start Up
log_file="mainnet_start_up.log"
echo "Starting Mainnet Setup..." | tee $log_file

# Step 1: Network Setup
echo "Setting up the network..."
network_setup_scripts=(
    "server_setup.sh"
    "network_setup.sh"
    "connection_pool_setup.sh"
    "distributed_network_co_ordinator_setup.sh"
    "fault_tolerance_setup.sh"
    "peer_advertisement_setup.sh"
    "qos_manager_setup.sh"
    "routing_setup.sh"
    "sdn_controller_setup.sh"
    "setup_peer_connection.sh"
    "setup_peer_discovery.sh"
    "tls_sls_setup.sh"
    "topology_setup.sh"
    "firewall_setup.sh"
)
for script in "${network_setup_scripts[@]}"; do
    run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/network_setup_scripts/$script" $log_file
done

# Step 2: Ledger Setup
echo "Setting up the ledger..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/ledger_set_up.sh" $log_file

# Step 3: High Availability Setup
echo "Setting up high availability..."
high_availability_scripts=(
    "api_replication_script.sh"
    "cli_replication_script.sh"
    "data_backup_script.sh"
    "data_collection_setup.sh"
    "data_distribution_setup.sh"
    "data_replication_script.sh"
    "heartbeat_service_setup.sh"
    "node_monitoring_setup.sh"
    "node_synchronization_script.sh"
    "redundancy_setup.sh"
)
for script in "${high_availability_scripts[@]}"; do
    run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/high_availability_scripts/$script" $log_file
done

# Step 4: Coin Setup
echo "Setting up the coin..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/coin_set_up.sh" $log_file

# Step 5: Owner Wallet Setup
echo "Setting up the owner wallet..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/genesis_wallet_setups.sh" $log_file

# Step 6: Consensus Setup
echo "Setting up the consensus mechanism..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/consensus_setup.sh" $log_file

# Step 7: Genesis Block Setup
echo "Setting up the genesis block..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/blockchain_and_genesis_block_setup.sh" $log_file

# Step 8: Release Genesis Keys
echo "Releasing the genesis keys..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/authority_genesis_key_creation.sh" $log_file

# Step 9: Loan Pool Setup
echo "Setting up the loan pool..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/loanpool_initialization_setup.sh" $log_file

# Step 10: Charity Pool Setup
echo "Setting up the charity pool..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/charity_pool_setup.sh" $log_file

# Step 11: Virtual Machine Setup
echo "Setting up the virtual machine..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/virtual_machine_setup.sh" $log_file

# Step 12: Storage Setup
echo "Setting up storage..."
run_script "/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/scripts/mainnet_start_up_scripts/storage_setup/storage_setup.sh" $log_file



echo "Mainnet setup completed."
