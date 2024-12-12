#!/bin/bash

# Script for Cross Consensus Switch Automation
# This script will start the consensus switch, optimize it, and check its consistency continuously in the mainnet.

# Logging setup
LOGFILE="/var/log/cross_consensus_switch_automation.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/consensus/switch"

# Function to start the consensus switching process
start_consensus_switching() {
    echo "$(date) - Starting the consensus switch process..." >> "$LOGFILE"
    while true; do
        # API call to start the consensus switching process
        response=$(curl -s -X POST "$BASE_URL/start")
        echo "$(date) - Start Consensus Switching Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Function to optimize the consensus switch
optimize_consensus_switch() {
    echo "$(date) - Optimizing the consensus switch..." >> "$LOGFILE"
    while true; do
        # API call to optimize the consensus switch process
        response=$(curl -s -X POST "$BASE_URL/optimize")
        echo "$(date) - Optimize Consensus Switch Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Function to check the consistency of the consensus switch
check_consensus_consistency() {
    echo "$(date) - Checking the consistency of the consensus switch..." >> "$LOGFILE"
    while true; do
        # API call to check the consistency of the consensus switch process
        response=$(curl -s -X POST "$BASE_URL/check-consistency")
        echo "$(date) - Check Consensus Consistency Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Main process
echo "$(date) - Initializing cross consensus switch automation..." >> "$LOGFILE"

# Start consensus switching process
start_consensus_switching &

# Optimize consensus switching process
optimize_consensus_switch &

# Check consensus switch consistency
check_consensus_consistency &

# Keep the script running
wait
