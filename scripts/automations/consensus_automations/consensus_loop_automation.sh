#!/bin/bash

# Script for Consensus Loop Automation
# This script will continuously manage the consensus loop, including starting/stopping execution,
# checking sub-block counts, and triggering block finalization.

# Logging setup
LOGFILE="/var/log/consensus_loop_automation.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/consensus"

# Function to start consensus execution
start_consensus_execution() {
    echo "$(date) - Starting consensus execution..." >> "$LOGFILE"
    while true; do
        # API call to start consensus execution
        response=$(curl -s -X POST "$BASE_URL/start-execution")
        echo "$(date) - Start Consensus Execution Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Function to stop consensus execution
stop_consensus_execution() {
    echo "$(date) - Stopping consensus execution..." >> "$LOGFILE"
    while true; do
        # API call to stop consensus execution
        response=$(curl -s -X POST "$BASE_URL/stop-execution")
        echo "$(date) - Stop Consensus Execution Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Function to check sub-block count
check_sub_block_count() {
    echo "$(date) - Checking sub-block count..." >> "$LOGFILE"
    while true; do
        # API call to check sub-block count
        response=$(curl -s "$BASE_URL/check-sub-block-count")
        echo "$(date) - Sub-block Count Check Response: $response" >> "$LOGFILE"

        # Logic to determine if block finalization is needed based on sub-block count
        if [[ "$response" -ge 1000 ]]; then
            echo "$(date) - Sub-block count exceeds threshold. Triggering block finalization..." >> "$LOGFILE"
            trigger_block_finalization
        else
            echo "$(date) - Sub-block count is below threshold." >> "$LOGFILE"
        fi
        # Continuous loop without sleep
    done
}

# Function to trigger block finalization
trigger_block_finalization() {
    echo "$(date) - Triggering block finalization..." >> "$LOGFILE"
    while true; do
        # API call to trigger block finalization
        response=$(curl -s -X POST "$BASE_URL/trigger-block-finalization")
        echo "$(date) - Block Finalization Trigger Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Main process
echo "$(date) - Initializing consensus loop automation..." >> "$LOGFILE"

# Start consensus execution in the background
start_consensus_execution &

# Check sub-block count and trigger block finalization if necessary
check_sub_block_count &

# Keep the script running
wait
