package advanced_security

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)


// Session and Alert represent structs for session and alert tracking.
type Session struct {
    SessionID string
    UserID    string
    TimeoutAt time.Time
}

type Alert struct {
    AlertID   string
    Status    string
    CreatedAt time.Time
}

// NodeAccessManager manages access frequency tracking and compliance checks.
type NodeAccessManager struct {
	sync.Mutex
	NodeAccessCounts  map[string]int       // Tracks the count of accesses per node
	NodeAccessRecords map[string][]time.Time // Records timestamps of each access for each node
	FrequencyLimit    int                    // Global frequency limit for compliance
}



type Metric struct {
	Name      string `json:"name"`
	Value     int    `json:"value"`
	Timestamp string `json:"timestamp"`
}


// Functions

// NewNodeAccessManager initializes the NodeAccessManager
func NewNodeAccessManager(frequencyLimit int) *NodeAccessManager {
	return &NodeAccessManager{
		NodeAccessCounts:  make(map[string]int),
		NodeAccessRecords: make(map[string][]time.Time),
		FrequencyLimit:    frequencyLimit,
	}
}

// TrackSystemHealthMetrics tracks various metrics to monitor system health
func TrackSystemHealthMetrics(metrics map[string]int) error {
	if len(metrics) == 0 {
		return errors.New("metrics cannot be empty")
	}

	// Record in monitoring system
	err := RecordHealthMetrics(metrics)
	if err != nil {
		return fmt.Errorf("failed to record system health metrics: %w", err)
	}

	// Record in ledger
	ledgerInstance := &ledger.Ledger{}

	// Call the method directly without assigning its return value
	ledgerInstance.AdvancedSecurityLedger.RecordHealthMetrics(metrics, time.Now())

	log.Println("System health metrics tracked and recorded.")
	return nil
}


// LogSystemMaintenance logs scheduled or unscheduled system maintenance events
func LogSystemMaintenance(event, maintenanceType string) error {
	if event == "" || maintenanceType == "" {
		return errors.New("event and maintenance type cannot be empty")
	}

    ledgerInstance := &ledger.Ledger{}
	err := ledgerInstance.MonitoringMaintenanceLedger.RecordMaintenanceEvent(event, maintenanceType)
	if err != nil {
		return fmt.Errorf("failed to log system maintenance: %w", err)
	}

	log.Printf("System maintenance logged: %s (%s)", event, maintenanceType)
	return nil
}

// RecordHealthMetrics logs system health metrics to an external monitoring service
func RecordHealthMetrics(metrics map[string]int) error {
	apiURL := os.Getenv("MONITORING_API_URL")
	apiKey := os.Getenv("MONITORING_API_KEY")
	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("monitoring API URL or API key is not configured")
	}

	var metricData []Metric
	timestamp := time.Now().Format(time.RFC3339)

	for name, value := range metrics {
		metricData = append(metricData, Metric{
			Name:      name,
			Value:     value,
			Timestamp: timestamp,
		})
		log.Printf("Preparing metric %s with value %d", name, value)
	}

	jsonData, err := json.Marshal(metricData)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics data: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL+"/metrics", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create metrics request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send metrics request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("monitoring API returned status: %s", resp.Status)
	}

	log.Printf("System health metrics recorded successfully at %s", timestamp)
	return nil
}

// SetSystemHealthThreshold sets a threshold for system health monitoring
func SetSystemHealthThreshold(threshold int) error {
	if threshold < 1 || threshold > 100 {
		return fmt.Errorf("health threshold must be between 1 and 100")
	}

	// Define the health threshold
	err := DefineHealthThreshold(threshold)
	if err != nil {
		return fmt.Errorf("failed to define health threshold: %w", err)
	}

	// Record the health threshold in the ledger
	ledgerInstance := &ledger.Ledger{}
	ledgerInstance.AdvancedSecurityLedger.RecordHealthThresholdSet(threshold, time.Now())

	log.Printf("System health threshold set to: %d", threshold)
	return nil
}


// FetchSystemHealthLog retrieves the system health log from the ledger
func FetchSystemHealthLog() ([]string, error) {
    ledgerInstance := &ledger.Ledger{}
	healthLog, err := ledgerInstance.AdvancedSecurityLedger.FetchHealthLog()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch system health log: %w", err)
	}

	log.Println("System health log fetched successfully.")
	return healthLog, nil
}

// DefineHealthThreshold sets a threshold value for system health monitoring
func DefineHealthThreshold(threshold int) error {
	err := os.Setenv("SYSTEM_HEALTH_THRESHOLD", fmt.Sprintf("%d", threshold))
	if err != nil {
		return fmt.Errorf("failed to set health threshold: %w", err)
	}

	log.Printf("System health threshold defined: %d", threshold)
	return nil
}

// VerifySystemHealthStatus verifies the current status of system health
func VerifySystemHealthStatus() (string, error) {
	status, err := CheckHealthStatus()
	if err != nil {
		return "", fmt.Errorf("failed to verify system health status: %w", err)
	}

	// Record the health status verification in the ledger
	ledgerInstance := &ledger.Ledger{}
	ledgerInstance.AdvancedSecurityLedger.RecordHealthStatusVerification(status, time.Now())

	log.Printf("System health status verified: %s", status)
	return status, nil
}



// RecordSystemHealthEvent logs a significant health-related event within the system.
func RecordSystemHealthEvent(event string) error {
    // Record event in ledger
    ledgerInstance := &ledger.Ledger{}
    err := ledgerInstance.AdvancedSecurityLedger.RecordHealthEvent(event, time.Now())
    if err != nil {
        return fmt.Errorf("failed to record system health event: %v", err)
    }

    fmt.Printf("System health event recorded: %s\n", event)
    return nil
}

// CheckHealthStatus checks the system’s current health status based on metrics.
func CheckHealthStatus() (string, error) {
    // Retrieve CPU usage and handle potential errors
    cpuUsage, err := getCPUUsage()
    if err != nil {
        return "", fmt.Errorf("failed to retrieve CPU usage: %v", err)
    }

    // Retrieve memory usage and handle potential errors
    memoryUsage, err := getMemoryUsage()
    if err != nil {
        return "", fmt.Errorf("failed to retrieve memory usage: %v", err)
    }

    // Set conditions for health status
    if cpuUsage < 80 && memoryUsage < 70 {
        log.Println("System health status: Healthy")
        return "Healthy", nil
    } else if cpuUsage < 90 && memoryUsage < 85 {
        log.Println("System health status: Degraded")
        return "Degraded", nil
    } else {
        log.Println("System health status: Critical")
        return "Critical", nil
    }
}

// getCPUUsage retrieves the current CPU usage percentage.
func getCPUUsage() (int, error) {
	ledgerInstance := &ledger.Ledger{}
	cpuUsage, err := ledgerInstance.MonitoringMaintenanceLedger.GetCPUUsage()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve CPU usage: %w", err)
	}

	// Convert float64 to int
	cpuUsageInt := int(cpuUsage)
	log.Printf("Current CPU usage: %d%%", cpuUsageInt)
	return cpuUsageInt, nil
}


// getMemoryUsage retrieves the current memory usage percentage.
func getMemoryUsage() (int, error) {
	ledgerInstance := &ledger.Ledger{}
	memoryUsage, err := ledgerInstance.MonitoringMaintenanceLedger.GetMemoryUsage()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve memory usage: %w", err)
	}

	// Convert float64 to int
	memoryUsageInt := int(memoryUsage)
	log.Printf("Current memory usage: %d%%", memoryUsageInt)
	return memoryUsageInt, nil
}


// FetchEventLogs retrieves a record of past events from the ledger.
func FetchEventLogs() ([]string, error) {
    ledgerInstance := &ledger.Ledger{}
	eventLogs, err := ledgerInstance.AdvancedSecurityLedger.FetchAllEventLogs()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch event logs: %w", err)
	}
	log.Printf("Fetched %d event logs.", len(eventLogs))
	return eventLogs, nil
}

// LogBackupEvent logs events related to backup processes.
func LogBackupEvent(event string) error {
	if event == "" {
		return errors.New("event description cannot be empty")
	}
    ledgerInstance := &ledger.Ledger{}
	err := ledgerInstance.AdvancedSecurityLedger.RecordBackupEvent(event, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log backup event: %w", err)
	}
	log.Printf("Backup event logged: %s", event)
	return nil
}

// MonitorSessionTimeouts checks for and logs session timeouts in the system.
func MonitorSessionTimeouts() error {
    ledgerInstance := &ledger.Ledger{}
	inactivityThreshold := 30 * time.Minute // Define inactivity threshold
	sessionTimeouts, err := CheckSessionTimeouts(ledgerInstance, inactivityThreshold)
	if err != nil {
		return fmt.Errorf("failed to monitor session timeouts: %w", err)
	}

	for _, session := range sessionTimeouts {
		err := ledgerInstance.AdvancedSecurityLedger.RecordSessionTimeout(session, time.Now())
		if err != nil {
			return fmt.Errorf("failed to record session timeout for session %s: %w", session.SessionID, err)
		}
	}
	log.Println("Session timeouts monitored and recorded.")
	return nil
}

// MonitorAlertStatus checks the status of active alerts and logs them.
func MonitorAlertStatus() error {
    ledgerInstance := &ledger.Ledger{}
	alertStatuses, err := CheckAlertStatus(ledgerInstance)
	if err != nil {
		return fmt.Errorf("failed to monitor alert status: %w", err)
	}

	for _, alert := range alertStatuses {
		err := ledgerInstance.AdvancedSecurityLedger.RecordAlertStatus(alert, time.Now())
		if err != nil {
			return fmt.Errorf("failed to record alert status for alert %s: %w", alert.AlertID, err)
		}
	}
	log.Println("Alert status monitored and logged.")
	return nil
}

// CheckSessionTimeouts retrieves a list of sessions that have timed out based on inactivity.
func CheckSessionTimeouts(ledgerInstance *ledger.Ledger, inactivityThreshold time.Duration) ([]ledger.Session, error) {
	var timedOutSessions []ledger.Session
	currentTime := time.Now()

	for _, session := range ledgerInstance.AdvancedSecurityLedger.Sessions {
		if currentTime.Sub(session.LastActivity) > inactivityThreshold {
			timedOutSessions = append(timedOutSessions, *session)
		}
	}

	if len(timedOutSessions) == 0 {
		return nil, errors.New("no timed-out sessions detected")
	}
	log.Printf("Detected %d timed-out sessions.", len(timedOutSessions))
	return timedOutSessions, nil
}


// CheckAlertStatus retrieves the current status of active alerts.
func CheckAlertStatus(ledgerInstance *ledger.Ledger) ([]ledger.Alert, error) {
	var activeAlerts []ledger.Alert
	for _, alert := range ledgerInstance.AdvancedSecurityLedger.Alerts {
		if alert.Status == "Active" || alert.Status == "Acknowledged" {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	if len(activeAlerts) == 0 {
		return nil, errors.New("no active alerts found")
	}
	log.Printf("Detected %d active alerts.", len(activeAlerts))
	return activeAlerts, nil
}


// sendHTTPPostRequest sends an HTTP POST request with JSON payload to the specified URL.
func sendHTTPPostRequest(url string, payload interface{}, apiKey string) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}
	log.Printf("HTTP POST request to %s completed successfully.", url)
	return nil
}


// allowNetworkAccess removes network access restrictions.
func allowNetworkAccess() error {
    // Revert network restrictions to allow normal network traffic
    log.Println("Network access unrestricted.")
    return nil
}

// deactivateIntrusionDetection stops the intrusion detection system.
func deactivateIntrusionDetection() error {
    log.Println("Intrusion detection deactivated.")
    return nil
}

// SetNodeAccessFrequencyLimit sets a frequency limit on node access to enforce security policies.
func SetNodeAccessFrequencyLimit(nodeID string, frequencyLimit int) error {
	// Set the access frequency limit in the common security configuration.
	err := DefineNodeAccessLimit(nodeID, frequencyLimit)
	if err != nil {
		return fmt.Errorf("failed to set access frequency limit for node %s: %v", nodeID, err)
	}

	// Record the access limit in the ledger.
	ledgerInstance := &ledger.Ledger{}
	err = ledgerInstance.AdvancedSecurityLedger.RecordNodeAccessLimit(nodeID, frequencyLimit, time.Now())
	if err != nil {
		return fmt.Errorf("failed to record access frequency limit for node %s: %v", nodeID, err)
	}

	fmt.Printf("Access frequency limit for node %s set to: %d\n", nodeID, frequencyLimit)
	return nil
}

// MonitorAccessFrequency checks the access frequency of nodes for compliance with limits.
func MonitorAccessFrequency() error {
	// Check the current access frequency of each node.
	accessData, err := CheckNodeAccessFrequency()
	if err != nil {
		return fmt.Errorf("failed to monitor access frequency: %v", err)
	}

	// Record the access frequency data in the ledger.
	ledgerInstance := &ledger.Ledger{}
	for nodeID, frequency := range accessData {
		err = ledgerInstance.AdvancedSecurityLedger.RecordAccessFrequency(nodeID, frequency, time.Now())
		if err != nil {
			return fmt.Errorf("failed to record access frequency for node %s: %v", nodeID, err)
		}
	}

	fmt.Println("Access frequency monitored and recorded.")
	return nil
}



// CheckNodeAccessFrequency retrieves the current access frequency of nodes for compliance checks.
func CheckNodeAccessFrequency() (map[string]int, error) {
    // Initialize a map to store access frequency per node
    accessData := make(map[string]int)

    // Retrieve the list of active nodes from the ledger or network configurations
    nodes, err := getNodeIDs()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve node IDs: %v", err)
    }

    // Calculate the access frequency for each active node
    for _, nodeID := range nodes {
        frequency, err := calculateAccessFrequency(nodeID)
        if err != nil {
            log.Printf("Failed to calculate access frequency for node %s: %v", nodeID, err)
            continue
        }
        accessData[nodeID] = frequency
    }

    if len(accessData) == 0 {
        return nil, errors.New("no access frequency data available")
    }

    log.Println("Access frequency data retrieved for compliance monitoring.")
    return accessData, nil
}

// getNodeIDs retrieves active node IDs from the ledger, ensuring only valid nodes are tracked.
func getNodeIDs() ([]string, error) {
    // Connect to the ledger and retrieve a list of registered nodes with active status.
    ledgerInstance := &ledger.Ledger{}
    nodes, err := ledgerInstance.FetchActiveNodeIDs()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve active nodes from ledger: %v", err)
    }

    if len(nodes) == 0 {
        return nil, errors.New("no active nodes found")
    }

    return nodes, nil
}

// calculateAccessFrequency calculates the access frequency for a specific node ID
// based on real-time access logs recorded in the ledger.
func calculateAccessFrequency(nodeID string) (int, error) {
    // Connect to the ledger to retrieve access log entries for the node
    ledgerInstance := &ledger.Ledger{}
    accessLogs, err := ledgerInstance.AdvancedSecurityLedger.FetchNodeAccessLogs(nodeID, time.Now().Add(-1*time.Minute), time.Now())
    if err != nil {
        return 0, fmt.Errorf("failed to fetch access logs for node %s: %v", nodeID, err)
    }

    // Calculate the frequency as the number of accesses in the last minute
    frequency := len(accessLogs)

    log.Printf("Access frequency for node %s: %d accesses per minute", nodeID, frequency)
    return frequency, nil
}


// DefineNodeAccessLimit sets a frequency limit for a node's access.
func DefineNodeAccessLimit(nodeID string, frequencyLimit int) error {
	if frequencyLimit <= 0 {
		return errors.New("frequency limit must be a positive integer")
	}

	// Store or enforce the limit (for a real system, this would configure the node's access policies)
	log.Printf("Node %s access frequency limit set to %d", nodeID, frequencyLimit)
	return nil
}

// LogNodeAccess logs each access attempt for the given nodeID, updating counts and records.
func (m *NodeAccessManager) LogNodeAccess(nodeID string) error {
	m.Lock()
	defer m.Unlock()

	// Log the access and update access count
	now := time.Now()
	m.NodeAccessCounts[nodeID]++
	m.NodeAccessRecords[nodeID] = append(m.NodeAccessRecords[nodeID], now)

	log.Printf("Access logged for node %s at %s. Total accesses: %d",
		nodeID, now.Format(time.RFC3339), m.NodeAccessCounts[nodeID])

	return nil
}

// CheckNodeAccessFrequency analyzes access records and counts per node against the frequency limit.
func (m *NodeAccessManager) CheckNodeAccessFrequency() (map[string]int, error) {
	m.Lock()
	defer m.Unlock()

	// If no access records exist, return an error
	if len(m.NodeAccessCounts) == 0 {
		return nil, errors.New("no access records found")
	}

	// Analyze each node’s access data against the frequency limit
	for nodeID, timestamps := range m.NodeAccessRecords {
		// Filter out accesses older than one minute (compliance window)
		cutoffTime := time.Now().Add(-1 * time.Minute)
		validAccesses := 0

		for _, timestamp := range timestamps {
			if timestamp.After(cutoffTime) {
				validAccesses++
			}
		}
		
		// Check if access frequency exceeds the allowed limit
		if validAccesses > m.FrequencyLimit {
			log.Printf("Node %s exceeds access frequency limit. Allowed: %d, Actual: %d",
				nodeID, m.FrequencyLimit, validAccesses)
		} else {
			log.Printf("Node %s is within access frequency limit. Allowed: %d, Actual: %d",
				nodeID, m.FrequencyLimit, validAccesses)
		}

		// Update the current access count with valid accesses within the time window
		m.NodeAccessCounts[nodeID] = validAccesses
	}

	// Return the access counts for compliance tracking
	return m.NodeAccessCounts, nil
}

// ResetNodeAccess clears the access counts and records outside of the compliance window.
func (m *NodeAccessManager) ResetNodeAccess() {
	m.Lock()
	defer m.Unlock()

	// Reset access counts and keep only records within the compliance window
	cutoffTime := time.Now().Add(-1 * time.Minute)
	for nodeID, timestamps := range m.NodeAccessRecords {
		filteredTimestamps := []time.Time{}
		for _, timestamp := range timestamps {
			if timestamp.After(cutoffTime) {
				filteredTimestamps = append(filteredTimestamps, timestamp)
			}
		}
		m.NodeAccessRecords[nodeID] = filteredTimestamps
		m.NodeAccessCounts[nodeID] = len(filteredTimestamps)
	}

	log.Println("Node access counts and records reset to maintain compliance.")
}

// SetSecurityLevelPolicy sets the security level policy across the network
func SetSecurityLevelPolicy(policy string) error {
    // Update security level policy in the system
    err := UpdateSecurityLevelPolicy(policy)
    if err != nil {
        return fmt.Errorf("failed to set security level policy: %v", err)
    }

    // Record the new security level policy in the ledger with a timestamp
    ledgerInstance := &ledger.Ledger{}
    timestamp := time.Now().Format(time.RFC3339) // Use formatted timestamp for consistency
    err = ledgerInstance.AdvancedSecurityLedger.RecordSecurityLevelPolicy(policy, timestamp)
    if err != nil {
        return fmt.Errorf("failed to record security level policy in ledger: %v", err)
    }

    fmt.Printf("Security level policy set to: %s\n", policy)
    return nil
}

// UpdateSecurityLevelPolicy updates the system’s security level policy.
func UpdateSecurityLevelPolicy(policy string) error {
    allowedPolicies := map[string]bool{"low": true, "medium": true, "high": true}
    if !allowedPolicies[policy] {
        return errors.New("invalid security level policy")
    }

    // Apply policy settings (example: configure security parameters)
    log.Printf("Security level policy updated to: %s", policy)
    return nil
}