#!/bin/bash

# Script for Dynamic Consensus Switch Automation
# This script will continuously handle starting the consensus switch, allocating validators, and logging contributions
# in the mainnet. The script runs in a loop to ensure continuous execution and real-time participation in data replication.

# Logging setup
LOGFILE="/var/log/dynamic_switch_automation.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/consensus/switch"

# Function to start execution monitoring
start_execution_monitoring() {
    echo "$(date) - Starting execution monitoring for the consensus switch..." >> "$LOGFILE"
    while true; do
        # API call to start execution monitoring
        response=$(curl -s -X POST "$BASE_URL/start")
        echo "$(date) - Start Execution Monitoring Response: $response" >> "$LOGFILE"
    done
}

# Function to allocate validators
allocate_validators() {
    echo "$(date) - Allocating validators for the consensus switch..." >> "$LOGFILE"
    while true; do
        # API call to allocate validators
        response=$(curl -s -X POST "$BASE_URL/allocate")
        echo "$(date) - Allocate Validators Response: $response" >> "$LOGFILE"
    done
}

# Function to log validator contributions
log_validator_contributions() {
    echo "$(date) - Logging validator contributions..." >> "$LOGFILE"
    while true; do
        # API call to log validator contributions
        response=$(curl -s -X POST "$BASE_URL/log-contributions")
        echo "$(date) - Log Validator Contributions Response: $response" >> "$LOGFILE"
    done
}

# Main process
echo "$(date) - Initializing dynamic consensus switch automation..." >> "$LOGFILE"

# Start execution monitoring
start_execution_monitoring &

# Allocate validators
allocate_validators &

# Log validator contributions
log_validator_contributions &

# Keep the script running
wait
