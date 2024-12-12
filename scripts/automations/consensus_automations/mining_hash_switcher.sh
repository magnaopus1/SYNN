#!/bin/bash

# Script for Mining Hash Switcher Automation
# This script monitors and applies hash switching for the mining process.
# It continuously runs in the mainnet and ensures optimal hash switching.

# Logging setup
LOGFILE="/var/log/mining_hash_switcher.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/hash_switch"

# Function to start hash switch monitoring
start_hash_switch_monitoring() {
    echo "$(date) - Starting hash switch monitoring..." >> "$LOGFILE"
    while true; do
        # API call to start hash switch monitoring
        response=$(curl -s -X POST "$BASE_URL/start")
        echo "$(date) - Start Hash Switch Monitoring Response: $response" >> "$LOGFILE"
    done
}

# Function to manually switch the hash
manual_hash_switch() {
    echo "$(date) - Triggering manual hash switch..." >> "$LOGFILE"
    while true; do
        # API call to manually switch the hash
        response=$(curl -s -X POST "$BASE_URL/manual")
        echo "$(date) - Manual Hash Switch Response: $response" >> "$LOGFILE"
    done
}

# Function to apply the new hash
apply_hash() {
    echo "$(date) - Applying new hash..." >> "$LOGFILE"
    while true; do
        # API call to apply the new hash
        response=$(curl -s -X POST "$BASE_URL/apply")
        echo "$(date) - Apply Hash Response: $response" >> "$LOGFILE"
    done
}

# Main process
echo "$(date) - Initializing mining hash switcher automation..." >> "$LOGFILE"

# Start hash switch monitoring
start_hash_switch_monitoring &

# Manually trigger hash switch
manual_hash_switch &

# Apply the new hash
apply_hash &

# Keep the script running
wait
