#!/bin/bash

# Script for Health Reporting Automation
# This script continuously monitors and triggers health reporting for the network.

# Logging setup
LOGFILE="/var/log/health_reporting.log"

# Base URL for the API
BASE_URL="http://localhost:8080/api/automations/health-report"

# Function to start health reporting
start_health_reporting() {
    echo "$(date) - Starting health reporting..." >> "$LOGFILE"
    while true; do
        # API call to start health reporting
        response=$(curl -s -X POST "$BASE_URL/start")
        echo "$(date) - Start Health Reporting Response: $response" >> "$LOGFILE"
    done
}

# Function to trigger health report generation
trigger_health_report() {
    echo "$(date) - Triggering health report..." >> "$LOGFILE"
    while true; do
        # API call to trigger health report
        response=$(curl -s -X POST "$BASE_URL/trigger")
        echo "$(date) - Trigger Health Report Response: $response" >> "$LOGFILE"
    done
}

# Main process
echo "$(date) - Initializing health reporting automation..." >> "$LOGFILE"

# Start health reporting
start_health_reporting &

# Trigger health report generation
trigger_health_report &

# Keep the script running
wait
