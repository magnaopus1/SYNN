#!/bin/bash

# Script for consensus latency monitoring automation.
# This script will continuously monitor latency, check latency levels, and trigger optimization as needed.

# Logging setup
LOGFILE="/var/log/consensus_latency_monitoring.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/latency"

# Function to start latency monitoring
start_latency_monitoring() {
    echo "$(date) - Starting latency monitoring..." >> "$LOGFILE"
    while true; do
        # API call to start latency monitoring
        response=$(curl -s -X POST "$BASE_URL/start-monitoring")
        echo "$(date) - Start Latency Monitoring Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Function to check the current latency
check_latency() {
    echo "$(date) - Checking consensus latency..." >> "$LOGFILE"
    while true; do
        # API call to check latency
        response=$(curl -s "$BASE_URL/check")
        echo "$(date) - Latency Check Response: $response" >> "$LOGFILE"

        # If latency is above a threshold or an issue is detected, trigger optimization
        if [[ "$response" != "OK" ]]; then
            echo "$(date) - High latency detected, initiating optimization..." >> "$LOGFILE"
            trigger_latency_optimization
        else
            echo "$(date) - Latency levels are within acceptable range." >> "$LOGFILE"
        fi
        # Continuous loop without sleep
    done
}

# Function to trigger latency optimization
trigger_latency_optimization() {
    echo "$(date) - Triggering latency optimization..." >> "$LOGFILE"
    while true; do
        # API call to trigger latency optimization
        response=$(curl -s -X POST "$BASE_URL/trigger-optimization")
        echo "$(date) - Latency Optimization Trigger Response: $response" >> "$LOGFILE"
        # Continuous loop without sleep
    done
}

# Main process
echo "$(date) - Initializing consensus latency monitoring automation..." >> "$LOGFILE"
start_latency_monitoring &
check_latency &

# Keep the script running
wait
