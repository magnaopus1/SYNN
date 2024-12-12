#!/bin/bash

# Script for Mining Difficulty Automation
# This script continuously monitors and manages mining difficulty synchronization.

# Logging setup
LOGFILE="/var/log/mining_difficulty.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/difficulty"

# Function to start mining difficulty synchronization
start_difficulty_sync() {
    echo "$(date) - Starting mining difficulty synchronization..." >> "$LOGFILE"
    while true; do
        # API call to start difficulty synchronization
        response=$(curl -s -X POST "$BASE_URL/start")
        echo "$(date) - Start Difficulty Sync Response: $response" >> "$LOGFILE"
    done
}

# Function to request mining difficulty change
request_difficulty_change() {
    echo "$(date) - Requesting mining difficulty change..." >> "$LOGFILE"
    while true; do
        # API call to request a difficulty change
        response=$(curl -s -X POST "$BASE_URL/request-change")
        echo "$(date) - Request Difficulty Change Response: $response" >> "$LOGFILE"
    done
}

# Main process
echo "$(date) - Initializing mining difficulty automation..." >> "$LOGFILE"

# Start mining difficulty synchronization
start_difficulty_sync &

# Request mining difficulty change
request_difficulty_change &

# Keep the script running
wait
