#!/bin/bash

# Script to handle consensus fallback automation by interacting with the API endpoints
# The script will continuously run, invoking the necessary APIs

# Logging setup
LOGFILE="/var/log/consensus_fallback_automation.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/consensus_fallback"

# Function to start consensus monitoring
start_consensus_monitoring() {
    echo "$(date) - Starting consensus monitoring..." >> "$LOGFILE"
    while true; do
        # API call to start monitoring
        response=$(curl -s -X POST "$BASE_URL/start")
        echo "$(date) - Start Consensus Monitoring Response: $response" >> "$LOGFILE"
        # Continuous looping without sleep
    done
}

# Function to trigger consensus fallback if necessary
trigger_fallback() {
    echo "$(date) - Triggering consensus fallback..." >> "$LOGFILE"
    while true; do
        # API call to trigger fallback
        response=$(curl -s -X POST "$BASE_URL/trigger")
        echo "$(date) - Trigger Fallback Response: $response" >> "$LOGFILE"
        # Continuous looping without sleep
    done
}

# Function to check the health of the consensus process
health_check() {
    echo "$(date) - Running consensus health check..." >> "$LOGFILE"
    while true; do
        # API call to check health
        response=$(curl -s "$BASE_URL/healthcheck")
        echo "$(date) - Health Check Response: $response" >> "$LOGFILE"

        if [[ "$response" != "OK" ]]; then
            echo "$(date) - Consensus Health Check Failed, triggering fallback..." >> "$LOGFILE"
            trigger_fallback
        else
            echo "$(date) - Consensus Health OK" >> "$LOGFILE"
        fi
        # Continuous looping without sleep
    done
}

# Main process
echo "$(date) - Initializing consensus fallback automation..." >> "$LOGFILE"
start_consensus_monitoring &
health_check &

# Keep the script running
wait
