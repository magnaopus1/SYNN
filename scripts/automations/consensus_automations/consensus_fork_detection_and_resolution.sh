#!/bin/bash

# Script for consensus fork detection and resolution automation.
# This script will continuously run, interacting with the API endpoints to detect, resolve, and monitor forks.

# Logging setup
LOGFILE="/var/log/consensus_fork_detection_and_resolution.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/consensus_fork"

# Function to start fork detection
start_fork_detection() {
    echo "$(date) - Starting fork detection..." >> "$LOGFILE"
    while true; do
        # API call to start fork detection
        response=$(curl -s -X POST "$BASE_URL/start")
        echo "$(date) - Start Fork Detection Response: $response" >> "$LOGFILE"
        # Continuous looping without sleep
    done
}

# Function to trigger fork resolution
trigger_fork_resolution() {
    echo "$(date) - Triggering fork resolution..." >> "$LOGFILE"
    while true; do
        # API call to resolve the fork
        response=$(curl -s -X POST "$BASE_URL/resolve")
        echo "$(date) - Trigger Fork Resolution Response: $response" >> "$LOGFILE"
        # Continuous looping without sleep
    done
}

# Function to check the current fork status
check_fork_status() {
    echo "$(date) - Checking fork status..." >> "$LOGFILE"
    while true; do
        # API call to check the fork status
        response=$(curl -s "$BASE_URL/check")
        echo "$(date) - Fork Status Check Response: $response" >> "$LOGFILE"

        # If fork is detected, trigger resolution
        if [[ "$response" != "OK" ]]; then
            echo "$(date) - Fork detected, initiating resolution..." >> "$LOGFILE"
            trigger_fork_resolution
        else
            echo "$(date) - No fork detected, system stable." >> "$LOGFILE"
        fi
        # Continuous looping without sleep
    done
}

# Main process
echo "$(date) - Initializing consensus fork detection and resolution automation..." >> "$LOGFILE"
start_fork_detection &
check_fork_status &

# Keep the script running
wait
