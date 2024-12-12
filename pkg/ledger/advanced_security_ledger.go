package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
)


// RecordAlertPolicySet records a new alert policy in the ledger with versioning
func (l *AdvancedSecurityLedger) RecordAlertPolicySet(policy string, timestamp time.Time) error {
	// Input validation
	if policy == "" {
		return fmt.Errorf("policy cannot be empty")
	}

	// Ensure timestamp is valid
	if timestamp.IsZero() {
		return fmt.Errorf("timestamp cannot be zero")
	}

	l.Lock()
	defer l.Unlock()

	// Determine the new policy version
	newVersion := 1
	if l.CurrentAlertPolicy != nil {
		newVersion = l.CurrentAlertPolicy.Version + 1
	}

	// Create and save the new policy
	newPolicy := AlertPolicy{
		Policy:  policy,
		Version: newVersion,
		SetAt:   timestamp,
	}

	l.AlertPolicies = append(l.AlertPolicies, newPolicy)
	l.CurrentAlertPolicy = &newPolicy

	log.Printf("[INFO] New alert policy set: '%s' (Version: %d) at %s", policy, newVersion, timestamp.Format(time.RFC3339))
	return nil
}


// UpdateCurrentAlertPolicy sets a new policy as the current alert policy
func (l *AdvancedSecurityLedger) UpdateCurrentAlertPolicy(policy string) error {
	// Input validation
	if policy == "" {
		return fmt.Errorf("policy cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	timestamp := time.Now()

	// Determine the new version based on the most recent policy
	newVersion := 1
	if l.CurrentAlertPolicy != nil {
		newVersion = l.CurrentAlertPolicy.Version + 1
	}

	// Create and set the new policy
	newPolicy := AlertPolicy{
		Policy:  policy,
		Version: newVersion,
		SetAt:   timestamp,
	}

	// Update the ledger
	l.AlertPolicies = append(l.AlertPolicies, newPolicy)
	l.CurrentAlertPolicy = &newPolicy

	log.Printf("[INFO] Updated current alert policy to version %d: '%s'", newVersion, policy)
	return nil
}


// GetCurrentAlertPolicy retrieves the latest alert policy details
func (l *AdvancedSecurityLedger) GetCurrentAlertPolicy() (*AlertPolicy, error) {
	l.Lock()
	defer l.Unlock()

	// Check if a policy is currently set
	if l.CurrentAlertPolicy == nil {
		err := errors.New("no alert policy is currently set")
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	log.Printf("[INFO] Retrieved current alert policy: '%s' (Version: %d)", l.CurrentAlertPolicy.Policy, l.CurrentAlertPolicy.Version)
	return l.CurrentAlertPolicy, nil
}



// Helper function to generate unique IDs for ledger entries based on an identifier (like sender or creator) and timestamp.
func (l *AdvancedSecurityLedger) generateUniqueID(identifier string) string {
	// Input validation
	if identifier == "" {
		log.Printf("[WARN] Empty identifier provided for unique ID generation")
		identifier = "default"
	}

	// Create input string with identifier and current timestamp
	input := fmt.Sprintf("%s-%d", identifier, time.Now().UnixNano())

	// Hash the input using SHA-256
	hash := sha256.Sum256([]byte(input))

	// Encode the hash as a hexadecimal string
	uniqueID := hex.EncodeToString(hash[:])

	log.Printf("[INFO] Generated unique ID: %s", uniqueID)
	return uniqueID
}

// GetEventMonitoringStatus retrieves the current event monitoring status
func (l *AdvancedSecurityLedger) GetEventMonitoringStatus() (string, error) {
	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Validate if the status is set
	if l.EventMonitoringStatus.Status == "" {
		err := errors.New("event monitoring status is undefined")
		log.Printf("[ERROR] %v", err)
		return "", err
	}

	log.Printf("[INFO] Event monitoring status retrieved: %s", l.EventMonitoringStatus.Status)
	return l.EventMonitoringStatus.Status, nil
}


// RecordAlertSent logs an alert sent for a security event
func (l *AdvancedSecurityLedger) RecordAlertSent(alertID, message string) error {
	// Input validation
	if alertID == "" || message == "" {
		return fmt.Errorf("alert ID and message must not be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Record the alert
	l.Alerts[alertID] = message
	log.Printf("[INFO] Alert recorded: ID=%s, Message='%s'", alertID, message)
	return nil
}


// RecordMitigationPlanSet stores a mitigation plan in the ledger
func (l *AdvancedSecurityLedger) RecordMitigationPlanSet(planID, description string) error {
	// Input validation
	if planID == "" || description == "" {
		return fmt.Errorf("plan ID and description must not be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Store the mitigation plan
	l.MitigationPlans[planID] = MitigationPlan{
		PlanID:        planID,
		Description:   description,
		Activated:     false,
		Effectiveness: "Pending",
		SetAt:         time.Now(),
	}

	log.Printf("[INFO] Mitigation plan recorded: ID=%s, Description='%s'", planID, description)
	return nil
}


// MitigationPlanExists checks if a mitigation plan with the specified ID exists
func (l *AdvancedSecurityLedger) MitigationPlanExists(planID string) (bool, error) {
	// Input validation
	if planID == "" {
		return false, fmt.Errorf("plan ID must not be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Check if the mitigation plan exists
	_, exists := l.MitigationPlans[planID]
	log.Printf("[INFO] Mitigation plan existence check: ID=%s, Exists=%v", planID, exists)
	return exists, nil
}


// ActivateMitigationPlan activates a mitigation plan by setting the activation time
func (l *AdvancedSecurityLedger) ActivateMitigationPlan(planID string) error {
	// Input validation
	if planID == "" {
		return fmt.Errorf("plan ID must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Check if the mitigation plan exists
	plan, exists := l.MitigationPlans[planID]
	if !exists {
		err := fmt.Errorf("mitigation plan %s does not exist and cannot be activated", planID)
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Activate the plan and update the timestamp
	plan.Activated = true
	plan.ActivatedAt = time.Now()
	l.MitigationPlans[planID] = plan

	log.Printf("[INFO] Mitigation plan %s activated at %s", planID, plan.ActivatedAt)
	return nil
}


// RecordMitigationActivated marks a mitigation plan as activated
func (l *AdvancedSecurityLedger) RecordMitigationActivated(planID string) error {
	// Input validation
	if planID == "" {
		return fmt.Errorf("plan ID must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Check if the mitigation plan exists
	plan, exists := l.MitigationPlans[planID]
	if !exists {
		err := fmt.Errorf("mitigation plan %s does not exist", planID)
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Mark the plan as activated
	plan.Activated = true
	plan.ActivatedAt = time.Now()
	l.MitigationPlans[planID] = plan

	log.Printf("[INFO] Mitigation plan %s marked as activated", planID)
	return nil
}


// GetMitigationMetrics retrieves the metrics associated with a specific mitigation plan
func (l *AdvancedSecurityLedger) GetMitigationMetrics(planID string) (MitigationMetrics, error) {
	// Input validation
	if planID == "" {
		return MitigationMetrics{}, fmt.Errorf("plan ID must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Retrieve the metrics for the mitigation plan
	metrics, exists := l.MetricsData[planID]
	if !exists {
		err := fmt.Errorf("metrics data not found for mitigation plan %s", planID)
		log.Printf("[ERROR] %v", err)
		return MitigationMetrics{}, err
	}

	log.Printf("[INFO] Metrics retrieved for mitigation plan %s: %+v", planID, metrics)
	return metrics, nil
}


// UpdateMitigationMetrics updates the metrics for a specific mitigation plan based on new data
func (l *AdvancedSecurityLedger) UpdateMitigationMetrics(planID string, newMetrics MitigationMetrics) error {
	// Input validation
	if planID == "" {
		return fmt.Errorf("plan ID must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Check if the mitigation plan exists
	if _, exists := l.MitigationPlans[planID]; !exists {
		err := fmt.Errorf("mitigation plan %s does not exist", planID)
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Update metrics with a new evaluation timestamp
	newMetrics.LastEvaluation = time.Now()
	l.MetricsData[planID] = newMetrics

	log.Printf("[INFO] Metrics updated for mitigation plan %s: %+v", planID, newMetrics)
	return nil
}

// DeactivateMitigationPlan deactivates a mitigation plan by updating its status
func (l *AdvancedSecurityLedger) DeactivateMitigationPlan(planID string) error {
	// Input validation
	if planID == "" {
		return fmt.Errorf("plan ID must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Retrieve the mitigation plan
	plan, exists := l.MitigationPlans[planID]
	if !exists {
		err := fmt.Errorf("mitigation plan %s not found", planID)
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Check if the plan is already inactive
	if !plan.Activated {
		err := fmt.Errorf("mitigation plan %s is not currently active", planID)
		log.Printf("[WARNING] %v", err)
		return err
	}

	// Deactivate the plan and update the timestamp
	now := time.Now()
	plan.Activated = false
	plan.DeactivatedAt = &now
	l.MitigationPlans[planID] = plan

	log.Printf("[INFO] Mitigation plan %s deactivated at %s", planID, now)
	return nil
}


// RecordMitigationDeactivated marks a mitigation plan as deactivated
func (l *AdvancedSecurityLedger) RecordMitigationDeactivated(planID string) error {
	// Input validation
	if planID == "" {
		return fmt.Errorf("plan ID must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Retrieve the mitigation plan
	plan, exists := l.MitigationPlans[planID]
	if !exists {
		err := fmt.Errorf("mitigation plan %s not found", planID)
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Mark the plan as deactivated
	plan.Activated = false
	l.MitigationPlans[planID] = plan

	log.Printf("[INFO] Mitigation plan %s marked as deactivated", planID)
	return nil
}


// RecordMitigationEffectiveness logs the effectiveness of a mitigation plan
func (l *AdvancedSecurityLedger) RecordMitigationEffectiveness(planID, effectiveness string) error {
	// Input validation
	if planID == "" {
		return fmt.Errorf("plan ID must not be empty")
	}
	if effectiveness == "" {
		return fmt.Errorf("effectiveness description must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Retrieve the mitigation plan
	plan, exists := l.MitigationPlans[planID]
	if !exists {
		err := fmt.Errorf("mitigation plan %s not found", planID)
		log.Printf("[ERROR] %v", err)
		return err
	}

	// Record the effectiveness
	plan.Effectiveness = effectiveness
	l.MitigationPlans[planID] = plan

	log.Printf("[INFO] Effectiveness of mitigation plan %s recorded as: %s", planID, effectiveness)
	return nil
}



// RecordEventMonitoringStatus logs the status of event monitoring
func (l *AdvancedSecurityLedger) RecordEventMonitoringStatus(eventID, status string) error {
	// Input validation
	if eventID == "" {
		return fmt.Errorf("event ID must not be empty")
	}
	if status == "" {
		return fmt.Errorf("status must not be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Record the event monitoring status
	l.Alerts[eventID] = status

	log.Printf("[INFO] Event monitoring status for %s recorded as: %s", eventID, status)
	return nil
}



// RecordNodeAccessLimit logs the access frequency limit set for a specific node.
func (l *AdvancedSecurityLedger) RecordNodeAccessLimit(nodeID string, frequencyLimit int, timestamp time.Time) error {
	// Input validation
	if nodeID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}
	if frequencyLimit <= 0 {
		return fmt.Errorf("frequency limit must be greater than zero")
	}

	l.Lock()
	defer l.Unlock()

	// Update the node access limits
	l.NodeAccessLimits[nodeID] = frequencyLimit

	// Log the operation
	log.Printf("[INFO] Access frequency limit recorded for node %s at %s: %d requests/min", nodeID, timestamp.Format(time.RFC3339), frequencyLimit)
	return nil
}


// RecordAccessFrequency logs the current access frequency of a specific node.
func (l *AdvancedSecurityLedger) RecordAccessFrequency(nodeID string, frequency int, timestamp time.Time) error {
	// Input validation
	if nodeID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}
	if frequency < 0 {
		return fmt.Errorf("frequency must be non-negative")
	}

	l.Lock()
	defer l.Unlock()

	// Update the access frequency
	l.AccessFrequencies[nodeID] = frequency

	// Log the operation
	log.Printf("[INFO] Access frequency recorded for node %s at %s: %d requests/min", nodeID, timestamp.Format(time.RFC3339), frequency)
	return nil
}


// RecordSecurityLevelPolicy records the security policy and timestamp in the ledger.
func (l *AdvancedSecurityLedger) RecordSecurityLevelPolicy(policy string, timestamp string) error {
	// Input validation
	if policy == "" {
		return fmt.Errorf("security policy cannot be empty")
	}
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	l.Lock()
	defer l.Unlock()

	// Append the new security policy
	l.SecurityPolicies = append(l.SecurityPolicies, SecurityPolicyRecord{
		Policy:    policy,
		Timestamp: timestamp,
	})

	// Log the operation
	log.Printf("[INFO] Security level policy recorded: %s at %s", policy, timestamp)
	return nil
}




// FetchNodeAccessLogs retrieves the access logs for a specific node within a given time range.
func (l *AdvancedSecurityLedger) FetchNodeAccessLogs(nodeID string, startTime, endTime time.Time) ([]AccessLog, error) {
	// Input validation
	if nodeID == "" {
		return nil, fmt.Errorf("node ID cannot be empty")
	}
	if startTime.After(endTime) {
		return nil, fmt.Errorf("start time must be before end time")
	}

	l.Lock()
	defer l.Unlock()

	var nodeAccessLogs []AccessLog
	for _, logEntry := range l.AccessLogs {
		// Filter logs by node ID and time range
		if logEntry.NodeID == nodeID && logEntry.Timestamp.After(startTime) && logEntry.Timestamp.Before(endTime) {
			nodeAccessLogs = append(nodeAccessLogs, logEntry)
		}
	}

	// Handle case where no logs are found
	if len(nodeAccessLogs) == 0 {
		return nil, fmt.Errorf("no access logs found for node %s within the specified time range", nodeID)
	}

	// Log the operation
	log.Printf("[INFO] Access logs for node %s retrieved for time range %s to %s", nodeID, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	return nodeAccessLogs, nil
}


// RecordFirmwareCheckStatus logs the firmware check status for a device
func (l *AdvancedSecurityLedger) RecordFirmwareCheckStatus(deviceID, status string) error {
	// Input validation
	if deviceID == "" {
		return fmt.Errorf("device ID cannot be empty")
	}
	if status == "" {
		return fmt.Errorf("status cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Update the firmware check status
	l.FirmwareCheckStatus[deviceID] = status

	// Log the operation
	log.Printf("[INFO] Firmware check status for device %s recorded as: %s", deviceID, status)
	return nil
}


// RecordConsensusAnomalyDetectionStatus logs the consensus anomaly detection status with a timestamp.
func (l *AdvancedSecurityLedger) RecordConsensusAnomalyDetectionStatus(status, timestamp string) error {
	// Input validation
	if status == "" {
		return fmt.Errorf("status cannot be empty")
	}
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	l.Lock()
	defer l.Unlock()

	// Update consensus anomaly detection status and timestamp
	l.ConsensusAnomalyDetectionStatus = status
	l.ConsensusAnomalyDetectionStatusTimestamp = timestamp

	// Log the operation
	log.Printf("[INFO] Consensus anomaly detection status set to %s at %s", status, timestamp)
	return nil
}


// FetchAllEventLogs retrieves all event logs from the ledger.
func (l *AdvancedSecurityLedger) FetchAllEventLogs() ([]string, error) {
	l.Lock()
	defer l.Unlock()

	// Check if there are any event logs
	if len(l.EventLogs) == 0 {
		return nil, fmt.Errorf("no event logs found")
	}

	// Format and retrieve logs
	var logs []string
	for _, logEntry := range l.EventLogs {
		formattedLog := fmt.Sprintf("[%s] %s", logEntry.Timestamp.Format(time.RFC3339), logEntry.Event)
		logs = append(logs, formattedLog)
	}

	// Log the operation
	log.Printf("[INFO] Retrieved %d event logs", len(logs))
	return logs, nil
}


// RecordBackupEvent logs a backup-related event in the ledger with a timestamp.
func (l *AdvancedSecurityLedger) RecordBackupEvent(event string, timestamp time.Time) error {
	// Input validation
	if event == "" {
		return fmt.Errorf("event description cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Create and append the new event log
	newEvent := EventLog{
		Event:     event,
		Timestamp: timestamp,
	}
	l.EventLogs = append(l.EventLogs, newEvent)

	// Log the operation
	log.Printf("[INFO] Backup event recorded: %s at %s", event, timestamp.Format(time.RFC3339))
	return nil
}

// RecordSessionTimeout logs a session timeout event in the ledger.
func (l *AdvancedSecurityLedger) RecordSessionTimeout(session Session, loggedAt time.Time) error {
	// Validate input
	if session.SessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	if session.UserID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Create and store the session timeout log
	logEntry := SessionTimeoutLog{
		SessionID: session.SessionID,
		UserID:    session.UserID,
		TimeoutAt: session.TimeoutAt,
		LoggedAt:  loggedAt,
	}
	l.SessionTimeoutLogs = append(l.SessionTimeoutLogs, logEntry)

	// Log the operation
	log.Printf("[INFO] Session timeout recorded for session ID: %s, user ID: %s at %s", session.SessionID, session.UserID, loggedAt.Format(time.RFC3339))
	return nil
}


// RecordIsolationIncident logs an isolation incident in the ledger.
func (l *AdvancedSecurityLedger) RecordIsolationIncident(incidentID string, timestamp time.Time) error {
	// Validate input
	if incidentID == "" {
		return fmt.Errorf("incident ID cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Create and store the isolation incident
	incident := IsolationIncident{
		IncidentID: incidentID,
		Timestamp:  timestamp,
	}
	l.IsolationIncidents = append(l.IsolationIncidents, incident)

	// Log the operation
	log.Printf("[INFO] Isolation incident logged for incident ID: %s at %s", incidentID, timestamp.Format(time.RFC3339))
	return nil
}


// RecordApplicationHardeningStatus logs the status of application hardening in the ledger.
func (l *AdvancedSecurityLedger) RecordApplicationHardeningStatus(status string, timestamp time.Time) error {
	// Validate input
	if status == "" {
		return fmt.Errorf("status cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Record the application hardening event
	event := ApplicationHardeningEvent{
		Status:    status,
		Timestamp: timestamp,
	}
	l.ApplicationHardeningEvents = append(l.ApplicationHardeningEvents, event)

	// Structured log for traceability
	log.Printf("[INFO] Application hardening status recorded: %s at %s", status, timestamp.Format(time.RFC3339))
	return nil
}



// RecordAlertStatus logs the status of an alert in the ledger.
func (l *AdvancedSecurityLedger) RecordAlertStatus(alert Alert, loggedAt time.Time) error {
	// Validate input
	if alert.AlertID == "" {
		return fmt.Errorf("alert ID cannot be empty")
	}
	if alert.Status == "" {
		return fmt.Errorf("alert status cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Create and store the alert status log
	logEntry := AlertStatusLog{
		AlertID:  alert.AlertID,
		Status:   alert.Status,
		LoggedAt: loggedAt,
	}
	l.AlertStatusLogs = append(l.AlertStatusLogs, logEntry)

	// Log the operation for traceability
	log.Printf("[INFO] Alert status recorded for alert ID: %s with status: %s at %s", alert.AlertID, alert.Status, loggedAt.Format(time.RFC3339))
	return nil
}


// RecordHealthEvent logs a significant health-related event with a timestamp.
func (l *AdvancedSecurityLedger) RecordHealthEvent(event string, timestamp time.Time) error {
	// Validate input
	if event == "" {
		return fmt.Errorf("health event description cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Create and store the health event log
	newEvent := HealthEvent{
		Event:     event,
		Timestamp: timestamp,
	}
	l.HealthEvents = append(l.HealthEvents, newEvent)

	// Log the operation for traceability
	log.Printf("[INFO] Health event recorded: %s at %s", event, timestamp.Format(time.RFC3339))
	return nil
}


// RecordHealthStatusVerification records a health status verification log in the ledger.
func (l *AdvancedSecurityLedger) RecordHealthStatusVerification(status string, timestamp time.Time) error {
	// Validate input
	if status == "" {
		return fmt.Errorf("health status cannot be empty")
	}

	l.Lock()
	defer l.Unlock()

	// Create and store the health status log
	logEntry := HealthStatusLog{
		Status:    status,
		Timestamp: timestamp,
	}
	l.HealthStatusLogs = append(l.HealthStatusLogs, logEntry)

	// Log the operation for traceability
	log.Printf("[INFO] Health status verified as %s at %s", status, timestamp.Format(time.RFC3339))
	return nil
}


// FetchHealthLog retrieves all system health log entries from the ledger.
func (l *AdvancedSecurityLedger) FetchHealthLog() ([]string, error) {
	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Validate if health logs exist
	if len(l.HealthLog) == 0 {
		return nil, fmt.Errorf("no health log entries found")
	}

	// Process and return health logs
	var healthLog []string
	for _, entry := range l.HealthLog {
		logEntry := fmt.Sprintf("[%s] %s", entry.Timestamp.Format(time.RFC3339), entry.Message)
		healthLog = append(healthLog, logEntry)
	}

	// Structured log for traceability
	log.Printf("[INFO] Retrieved %d health log entries.", len(healthLog))
	return healthLog, nil
}



// RecordAPIUsage logs API usage events
func (l *Ledger) RecordAPIUsage(apiID string, usageCount int) error {
	// Validate inputs
	if apiID == "" {
		return fmt.Errorf("API ID cannot be empty")
	}
	if usageCount < 0 {
		return fmt.Errorf("usage count cannot be negative")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Record API usage
	l.APILog[apiID] = APIUsageRecord{
		APIID:      apiID,
		UsageCount: usageCount,
		LastUsedAt: time.Now(),
	}

	// Log the operation for traceability
	log.Printf("[INFO] API usage recorded for %s with count %d at %s.", apiID, usageCount, time.Now().Format(time.RFC3339))
	return nil
}


// RecordAnomalyDetectionStatus logs the detection status of anomalies
func (l *AdvancedSecurityLedger) RecordAnomalyDetectionStatus(detectionID, status string) error {
	// Validate inputs
	if detectionID == "" {
		return fmt.Errorf("detection ID cannot be empty")
	}
	if status == "" {
		return fmt.Errorf("status cannot be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Append anomaly event
	event := AnomalyEvent{
		EventID:    detectionID,
		Type:       status,
		Severity:   1,
		DetectedAt: time.Now(),
	}
	l.AnomalyEvents[detectionID] = append(l.AnomalyEvents[detectionID], event)

	// Log the operation for traceability
	log.Printf("[INFO] Anomaly detection status for %s recorded as: %s at %s.", detectionID, status, time.Now().Format(time.RFC3339))
	return nil
}


// LogSuspiciousActivity records suspicious activity in the ledger
func (l *AdvancedSecurityLedger) LogSuspiciousActivity(activityID, description string) error {
	// Validate inputs
	if activityID == "" {
		return fmt.Errorf("activity ID cannot be empty")
	}
	if description == "" {
		return fmt.Errorf("description cannot be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Append suspicious activity record
	record := SuspiciousActivityRecord{
		ActivityID:  activityID,
		Description: description,
		DetectedAt:  time.Now(),
	}
	l.SuspiciousActivityLog = append(l.SuspiciousActivityLog, record)

	// Log the operation for traceability
	log.Printf("[INFO] Suspicious activity %s logged: %s at %s.", activityID, description, time.Now().Format(time.RFC3339))
	return nil
}



// RecordTrafficPattern logs normal traffic patterns
func (l *AdvancedSecurityLedger) RecordTrafficPattern(patternID, details string) error {
	// Input validation
	if patternID == "" {
		return fmt.Errorf("pattern ID cannot be empty")
	}
	if details == "" {
		return fmt.Errorf("details cannot be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Record the traffic pattern
	l.TrafficPatterns[patternID] = TrafficPattern{
		PatternID:  patternID,
		Details:    details,
		RecordedAt: time.Now(),
	}

	// Log the operation
	log.Printf("[INFO] Traffic pattern %s recorded with details: %s.", patternID, details)
	return nil
}


// RecordProtocolDeviation logs protocol deviations detected in the system
func (l *AdvancedSecurityLedger) RecordProtocolDeviation(deviationID, description string) error {
	// Input validation
	if deviationID == "" {
		return fmt.Errorf("deviation ID cannot be empty")
	}
	if description == "" {
		return fmt.Errorf("description cannot be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Record the protocol deviation
	l.ProtocolDeviations[deviationID] = ProtocolDeviation{
		DeviationID: deviationID,
		Description: description,
		DetectedAt:  time.Now(),
	}

	// Log the operation
	log.Printf("[INFO] Protocol deviation %s recorded with description: %s.", deviationID, description)
	return nil
}



// RecordAnomalyAlertLevelSet logs the anomaly alert level and timestamp in the ledger
func (l *AdvancedSecurityLedger) RecordAnomalyAlertLevelSet(level int, timestamp string) error {
	// Input validation
	if level < 0 {
		return fmt.Errorf("alert level cannot be negative")
	}
	if timestamp == "" {
		return fmt.Errorf("timestamp cannot be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Record the anomaly alert level
	l.AnomalyAlertLevel = level
	l.AlertTimestamp = timestamp

	// Log the operation
	log.Printf("[INFO] Anomaly alert level set to %d at %s.", level, timestamp)
	return nil
}


// RecordNetworkAlertSent logs an alert sent to the network
func (l *AdvancedSecurityLedger) RecordNetworkAlertSent(alertID, message string) error {
	// Input validation
	if alertID == "" {
		return fmt.Errorf("alert ID cannot be empty")
	}
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Record the network alert
	l.NetworkAlerts[alertID] = message

	// Log the operation
	log.Printf("[INFO] Network alert %s recorded with message: %s.", alertID, message)
	return nil
}

// LogAlertResponse logs the response to a specific alert
func (l *AdvancedSecurityLedger) LogAlertResponse(alertID, response string) error {
	// Input validation
	if alertID == "" {
		return fmt.Errorf("alert ID cannot be empty")
	}
	if response == "" {
		return fmt.Errorf("response cannot be empty")
	}

	// Ensure thread safety
	l.Lock()
	defer l.Unlock()

	// Record the alert response
	l.AlertResponses[alertID] = response

	// Log the operation
	log.Printf("[INFO] Response to alert %s logged: %s.", alertID, response)
	return nil
}


// RecordIntegrityViolation logs a detected integrity violation
func (l *AdvancedSecurityLedger) RecordIntegrityViolation(violationID, details string) error {
    // Input validation
    if violationID == "" {
        return fmt.Errorf("violation ID cannot be empty")
    }
    if details == "" {
        return fmt.Errorf("details cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the integrity violation
    l.IntegrityViolations[violationID] = details

    // Log the operation
    log.Printf("[INFO] Integrity violation %s recorded with details: %s.", violationID, details)
    return nil
}


// RecordThreatDetectionStatus sets the status of threat detection for a specific detection ID
func (l *AdvancedSecurityLedger) RecordThreatDetectionStatus(detectionID, status string) error {
    // Input validation
    if detectionID == "" {
        return fmt.Errorf("detection ID cannot be empty")
    }
    if status == "" {
        return fmt.Errorf("status cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the threat detection status
    l.ThreatDetectionStatus[detectionID] = ThreatDetection{
        DetectionID:   detectionID,
        Status:        status,
        LastCheckedAt: time.Now(),
    }

    // Log the operation
    log.Printf("[INFO] Threat detection status for %s set to %s.", detectionID, status)
    return nil
}


// RecordAnomalyEvent logs a detected anomaly event in the ledger
func (l *AdvancedSecurityLedger) RecordAnomalyEvent(event, timestamp string) error {
    // Input validation
    if event == "" {
        return fmt.Errorf("event cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the anomaly event
    anomaly := AnomalyEvent{
        Event:     event,
        Timestamp: timestamp,
    }
    l.AnomalyEvents = append(l.AnomalyEvents, anomaly)

    // Log the operation
    log.Printf("[INFO] Anomaly event recorded: %s at %s", event, timestamp)
    return nil
}


// RecordTrafficAnomaly logs a detected traffic anomaly in the ledger
func (l *AdvancedSecurityLedger) RecordTrafficAnomaly(details, timestamp string) error {
    // Input validation
    if details == "" {
        return fmt.Errorf("details cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the traffic anomaly
    anomaly := TrafficAnomaly{
        Details:   details,
        Timestamp: timestamp,
    }
    l.TrafficAnomalies = append(l.TrafficAnomalies, anomaly)

    // Log the operation
    log.Printf("[INFO] Traffic anomaly recorded: %s at %s", details, timestamp)
    return nil
}


// RecordThreatEvent logs a security threat event in the ledger
func (l *AdvancedSecurityLedger) RecordThreatEvent(event string, timestamp time.Time) error {
    // Input validation
    if event == "" {
        return fmt.Errorf("event cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the threat event
    l.ThreatEvents = append(l.ThreatEvents, ThreatEvent{
        Event:     event,
        Timestamp: timestamp,
    })

    // Log the operation
    log.Printf("[INFO] Threat event recorded: %s at %s", event, timestamp.Format(time.RFC3339))
    return nil
}


// RecordSecurityThresholdSet logs the security threshold in the ledger
func (l *AdvancedSecurityLedger) RecordSecurityThresholdSet(threshold int, timestamp string) error {
    // Validate input parameters
    if threshold <= 0 {
        return fmt.Errorf("security threshold must be greater than zero")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Update the security threshold
    l.SecurityThreshold = threshold
    l.ThresholdTimestamp = timestamp

    // Log the operation
    log.Printf("[INFO] Security alert threshold set to %d at %s", threshold, timestamp)
    return nil
}


// RecordSystemThreatLevel logs the system threat level in the ledger
func (l *AdvancedSecurityLedger) RecordSystemThreatLevel(level int, timestamp string) error {
    // Validate input parameters
    if level < 0 || level > 10 {
        return fmt.Errorf("threat level must be between 0 and 10")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the system threat level
    l.SystemThreatLevels = append(l.SystemThreatLevels, SystemThreat{
        Level:     level,
        Timestamp: timestamp,
    })

    // Log the operation
    log.Printf("[INFO] System threat level recorded: %d at %s", level, timestamp)
    return nil
}


// RecordThreatLevel logs a change in threat level in the ledger
func (l *AdvancedSecurityLedger) RecordThreatLevel(level int, timestamp string) error {
    // Validate input parameters
    if level < 0 || level > 10 {
        return fmt.Errorf("threat level must be between 0 and 10")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the threat level
    l.ThreatLevels = append(l.ThreatLevels, ThreatLevel{
        Level:     level,
        Timestamp: timestamp,
    })

    // Log the operation
    log.Printf("[INFO] Threat level recorded: %d at %s", level, timestamp)
    return nil
}


// RecordIncidentDeactivation logs the deactivation of an incident response in the ledger
func (l *AdvancedSecurityLedger) RecordIncidentDeactivation(incidentID, timestamp string) error {
    // Validate input parameters
    if incidentID == "" {
        return fmt.Errorf("incident ID cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the incident deactivation
    l.IncidentDeactivations = append(l.IncidentDeactivations, IncidentDeactivation{
        IncidentID: incidentID,
        Timestamp:  timestamp,
    })

    // Log the operation
    log.Printf("[INFO] Incident response deactivated for incident ID: %s at %s", incidentID, timestamp)
    return nil
}


// RecordIncidentEvent logs an event related to an incident response in the ledger
func (l *AdvancedSecurityLedger) RecordIncidentEvent(incidentID, event, timestamp string) error {
    // Validate input parameters
    if incidentID == "" {
        return fmt.Errorf("incident ID cannot be empty")
    }
    if event == "" {
        return fmt.Errorf("event cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    // Ensure thread safety
    l.Lock()
    defer l.Unlock()

    // Record the incident event
    l.IncidentEvents = append(l.IncidentEvents, IncidentEvent{
        IncidentID: incidentID,
        Event:      event,
        Timestamp:  timestamp,
    })

    // Log the operation
    log.Printf("[INFO] Incident event logged for incident ID: %s, Event: %s at %s", incidentID, event, timestamp)
    return nil
}


// RecordIncidentResolution logs the resolution status of an incident in the ledger
func (l *AdvancedSecurityLedger) RecordIncidentResolution(incidentID, resolutionStatus, timestamp string) error {
    // Validate inputs
    if incidentID == "" {
        return fmt.Errorf("incident ID cannot be empty")
    }
    if resolutionStatus == "" {
        return fmt.Errorf("resolution status cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    // Record the resolution status
    l.IncidentResolutions = append(l.IncidentResolutions, IncidentResolution{
        IncidentID:       incidentID,
        ResolutionStatus: resolutionStatus,
        Timestamp:        timestamp,
    })

    log.Printf("[INFO] Incident resolution status recorded: Incident ID %s, Status %s, Timestamp %s", incidentID, resolutionStatus, timestamp)
    return nil
}


// GetUnauthorizedAttempts retrieves the number of unauthorized attempts from the ledger
func (l *AdvancedSecurityLedger) GetUnauthorizedAttempts() (int, error) {
    l.Lock()
    defer l.Unlock()

    // Validate ledger state
    if l.UnauthorizedAttempts < 0 {
        return 0, fmt.Errorf("invalid number of unauthorized attempts in ledger: %d", l.UnauthorizedAttempts)
    }

    log.Printf("[INFO] Retrieved unauthorized attempts: %d", l.UnauthorizedAttempts)
    return l.UnauthorizedAttempts, nil
}


// GetSystemLoad retrieves the current system load from the ledger
func (l *AdvancedSecurityLedger) GetSystemLoad() (int, error) {
    l.Lock()
    defer l.Unlock()

    // Validate ledger state
    if l.SystemLoad < 0 || l.SystemLoad > 100 {
        return 0, fmt.Errorf("invalid system load in ledger: %d", l.SystemLoad)
    }

    log.Printf("[INFO] Retrieved system load: %d", l.SystemLoad)
    return l.SystemLoad, nil
}





// GetTrafficRate retrieves the current traffic rate from the ledger
func (l *AdvancedSecurityLedger) GetTrafficRate() (int, error) {
    l.Lock()
    defer l.Unlock()

    // Validate ledger state
    if l.TrafficRate < 0 {
        return 0, fmt.Errorf("invalid traffic rate in ledger: %d", l.TrafficRate)
    }

    log.Printf("[INFO] Retrieved traffic rate: %d", l.TrafficRate)
    return l.TrafficRate, nil
}


// RecordEscalationProtocolSet logs the escalation protocol in the ledger
func (l *AdvancedSecurityLedger) RecordEscalationProtocolSet(protocol, timestamp string) error {
    // Validate inputs
    if protocol == "" {
        return fmt.Errorf("protocol cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    // Record the escalation protocol
    l.EscalationProtocol = protocol
    l.EscalationTimestamp = timestamp

    log.Printf("[INFO] Escalation protocol set: %s at %s", protocol, timestamp)
    return nil
}


// GetDDoSScore retrieves the DDoS threat score from the ledger
func (l *AdvancedSecurityLedger) GetDDoSScore() (int, error) {
    l.Lock()
    defer l.Unlock()

    // Validate the score's range
    if l.DDoSScore < 0 || l.DDoSScore > 100 {
        return 0, fmt.Errorf("invalid DDoS score in ledger: %d", l.DDoSScore)
    }

    log.Printf("[INFO] Retrieved DDoS score: %d", l.DDoSScore)
    return l.DDoSScore, nil
}


// FetchAllIncidentReports retrieves all incident reports stored in the ledger
func (l *AdvancedSecurityLedger) FetchAllIncidentReports() ([]string, error) {
    l.Lock()
    defer l.Unlock()

    if len(l.IncidentReports) == 0 {
        return nil, fmt.Errorf("no incident reports found in the ledger")
    }

    // Collect and format reports for return
    reports := []string{}
    for _, report := range l.IncidentReports {
        reports = append(reports, fmt.Sprintf("Incident ID: %s, Details: %s", report.IncidentID, report.Details))
    }

    log.Printf("[INFO] Fetched %d incident reports.", len(reports))
    return reports, nil
}


// RecordRetentionPolicySet logs the retention policy for incident logs in the ledger
func (l *AdvancedSecurityLedger) RecordRetentionPolicySet(policyName, policy, timestamp string) error {
    // Validate inputs
    if policyName == "" {
        return fmt.Errorf("policy name cannot be empty")
    }
    if policy == "" {
        return fmt.Errorf("policy details cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    if l.RetentionPolicies == nil {
        l.RetentionPolicies = make(map[string]RetentionPolicy)
    }

    // Record the retention policy
    l.RetentionPolicies[policyName] = RetentionPolicy{
        PolicyName: policyName,
        Policy:     policy,
        Timestamp:  timestamp,
    }

    log.Printf("[INFO] Retention policy set for %s: %s at %s", policyName, policy, timestamp)
    return nil
}


// RecordIncidentResponseProtocol sets up an incident response protocol
func (l *AdvancedSecurityLedger) RecordIncidentResponseProtocol(protocolID, description string) error {
    // Validate inputs
    if protocolID == "" {
        return fmt.Errorf("protocol ID cannot be empty")
    }
    if description == "" {
        return fmt.Errorf("protocol description cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    // Record the incident response protocol
    l.IncidentProtocols[protocolID] = IncidentProtocol{
        ProtocolID:  protocolID,
        Description: description,
        Activated:   false,
    }

    log.Printf("[INFO] Incident response protocol %s set with description: %s.", protocolID, description)
    return nil
}


// RecordIncidentActivation logs the activation of an incident response
func (l *AdvancedSecurityLedger) RecordIncidentActivation(incidentID, timestamp string) error {
    // Validate inputs
    if incidentID == "" {
        return fmt.Errorf("incident ID cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    // Record the activation
    l.IncidentActivations = append(l.IncidentActivations, IncidentActivation{
        IncidentID: incidentID,
        Timestamp:  timestamp,
    })

    log.Printf("[INFO] Incident response activated for incident ID: %s at %s", incidentID, timestamp)
    return nil
}



// RecordIsolationEvent logs isolation events related to an isolation incident
func (l *AdvancedSecurityLedger) RecordIsolationEvent(incidentID, eventDetails string) error {
    if incidentID == "" {
        return fmt.Errorf("incident ID cannot be empty")
    }
    if eventDetails == "" {
        return fmt.Errorf("event details cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    if incident, exists := l.IsolationIncidents[incidentID]; exists {
        incident.Description += "; " + eventDetails
        l.IsolationIncidents[incidentID] = incident
        log.Printf("[INFO] Isolation event recorded for incident %s: %s", incidentID, eventDetails)
        return nil
    }

    return fmt.Errorf("incident ID %s not found in the ledger", incidentID)
}


// RecordHealthThresholdSet logs the system health threshold in the ledger with a timestamp.
func (l *AdvancedSecurityLedger) RecordHealthThresholdSet(threshold int, timestamp time.Time) error {
    if threshold <= 0 {
        return fmt.Errorf("health threshold must be positive")
    }

    l.Lock()
    defer l.Unlock()

    l.HealthThreshold = threshold
    l.HealthThresholdTimestamp = timestamp

    log.Printf("[INFO] Health threshold set to %d at %s", threshold, timestamp.Format(time.RFC3339))
    return nil
}


// RecordMaintenanceEvent logs a maintenance event with type and timestamp in the ledger.
func (l *Ledger) RecordMaintenanceEvent(event, maintenanceType string, timestamp time.Time) error {
    if event == "" {
        return fmt.Errorf("event description cannot be empty")
    }
    if maintenanceType == "" {
        return fmt.Errorf("maintenance type cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    maintenanceEvent := MaintenanceEvent{
        Event:          event,
        MaintenanceType: maintenanceType,
        Timestamp:      timestamp,
    }

    l.MaintenanceEvents = append(l.MaintenanceEvents, maintenanceEvent)

    log.Printf("[INFO] Maintenance event recorded: %s (%s) at %s", event, maintenanceType, timestamp.Format(time.RFC3339))
    return nil
}


// RecordHealthMetrics stores health metrics in the ledger with a timestamp.
func (l *AdvancedSecurityLedger) RecordHealthMetrics(metrics map[string]int, timestamp time.Time) error {
    if metrics == nil || len(metrics) == 0 {
        return fmt.Errorf("metrics cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    if l.HealthMetrics == nil {
        l.HealthMetrics = make(map[string]int)
    }

    for metric, value := range metrics {
        l.HealthMetrics[metric] = value
    }
    l.HealthMetricsTimestamp = timestamp

    log.Printf("[INFO] Health metrics recorded at %s: %v", timestamp.Format(time.RFC3339), metrics)
    return nil
}



// RecordAPILimitSet logs the API limit in the ledger with a timestamp
func (l *AdvancedSecurityLedger) RecordAPILimitSet(limit int, timestamp string) error {
    if limit <= 0 {
        return fmt.Errorf("API limit must be positive")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    l.APILimit = limit
    l.APILimitTimestamp = timestamp

    log.Printf("[INFO] API limit set to %d requests per second at %s", limit, timestamp)
    return nil
}


// RecordRateLimitingStatus logs the rate-limiting status in the ledger with a timestamp.
func (l *AdvancedSecurityLedger) RecordRateLimitingStatus(status, timestamp string) error {
    if status == "" {
        return fmt.Errorf("rate-limiting status cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    l.RateLimitingStatus = status
    l.RateLimitingStatusTimestamp = timestamp

    log.Printf("[INFO] Rate limiting status set to '%s' at %s", status, timestamp)
    return nil
}


// RecordTransferRateLimit logs the data transfer rate limit in the ledger with a timestamp.
func (l *AdvancedSecurityLedger) RecordTransferRateLimit(limit int, timestamp string) error {
    if limit <= 0 {
        return fmt.Errorf("data transfer rate limit must be greater than zero")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    l.TransferRateLimit = limit
    l.TransferRateLimitTimestamp = timestamp

    log.Printf("[INFO] Data transfer rate limit set to %d MB/s at %s", limit, timestamp)
    return nil
}




// RecordRateLimitPolicy logs the rate limit policy in the ledger
func (l *AdvancedSecurityLedger) RecordRateLimitPolicy(policy, timestamp string) error {
    if policy == "" {
        return fmt.Errorf("rate limit policy cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    l.RateLimitPolicy = policy
    l.RateLimitPolicyTimestamp = timestamp

    log.Printf("[INFO] Rate limit policy set to '%s' at %s", policy, timestamp)
    return nil
}


// RecordTransactionThreshold logs the transaction threshold in the ledger
func (l *AdvancedSecurityLedger) RecordTransactionThreshold(threshold int, timestamp string) error {
    if threshold <= 0 {
        return fmt.Errorf("transaction threshold must be greater than zero")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    l.TransactionThreshold = threshold
    l.TransactionThresholdTimestamp = timestamp

    log.Printf("[INFO] Transaction threshold set to %d at %s", threshold, timestamp)
    return nil
}


// RecordDetectedThreat logs a detected threat in the ledger
func (l *AdvancedSecurityLedger) RecordDetectedThreat(threatDetails, timestamp string) error {
    if threatDetails == "" {
        return fmt.Errorf("threat details cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    l.ThreatDetections = append(l.ThreatDetections, ThreatDetection{
        Details:   threatDetails,
        Timestamp: timestamp,
    })

    log.Printf("[INFO] Threat detection recorded: '%s' at %s", threatDetails, timestamp)
    return nil
}



// RecordIntrusionDetectionStatus logs the status of intrusion detection in the ledger
func (l *AdvancedSecurityLedger) RecordIntrusionDetectionStatus(status, timestamp string) error {
    if status == "" {
        return fmt.Errorf("intrusion detection status cannot be empty")
    }
    if timestamp == "" {
        return fmt.Errorf("timestamp cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    l.IntrusionDetectionStatus = status
    l.IntrusionDetectionTimestamp = timestamp

    log.Printf("[INFO] Intrusion detection status set to '%s' at %s", status, timestamp)
    return nil
}


// RecordSessionStart logs the start of a new session in the ledger
func (l *AdvancedSecurityLedger) RecordSessionStart(sessionID, userID, ipAddress string) error {
    if sessionID == "" || userID == "" || ipAddress == "" {
        return fmt.Errorf("session ID, user ID, and IP address cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    if _, exists := l.ActiveSessions[sessionID]; exists {
        return fmt.Errorf("session with ID '%s' already exists", sessionID)
    }

    newSession := &Session{
        SessionID:    sessionID,
        UserID:       userID,
        StartTime:    time.Now(),
        LastActivity: time.Now(),
        IPAddress:    ipAddress,
    }

    l.ActiveSessions[sessionID] = newSession

    log.Printf("[INFO] Session '%s' started for user '%s' from IP '%s'", sessionID, userID, ipAddress)
    return nil
}


// RecordSessionEnd logs the end of an existing session in the ledger.
func (l *AdvancedSecurityLedger) RecordSessionEnd(sessionID string) error {
    if sessionID == "" {
        return fmt.Errorf("session ID cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    session, exists := l.ActiveSessions[sessionID]
    if !exists {
        return fmt.Errorf("session with ID '%s' not found", sessionID)
    }

    session.LastActivity = time.Now()
    l.SessionLogs[sessionID] = session
    delete(l.ActiveSessions, sessionID)

    log.Printf("[INFO] Session '%s' ended", sessionID)
    return nil
}



// GetSessionLog retrieves the details of a session based on its session ID.
func (l *AdvancedSecurityLedger) GetSessionLog(sessionID string) (*Session, error) {
    if sessionID == "" {
        return nil, fmt.Errorf("session ID cannot be empty")
    }

    l.Lock()
    defer l.Unlock()

    session, exists := l.SessionLogs[sessionID]
    if !exists {
        return nil, fmt.Errorf("session log with ID '%s' not found", sessionID)
    }

    log.Printf("[INFO] Session log retrieved for session ID '%s'", sessionID)
    return session, nil
}


// IsMitigationPlanActive checks if a mitigation plan with the given ID is active.
func (asl *AdvancedSecurityLedger) IsMitigationPlanActive(planID string) (bool, error) {
    if planID == "" {
        return false, fmt.Errorf("plan ID cannot be empty")
    }

    asl.mu.Lock()
    defer asl.mu.Unlock()

    active, exists := asl.ActiveMitigationPlans[planID]
    if !exists {
        return false, fmt.Errorf("mitigation plan with ID '%s' not found", planID)
    }

    log.Printf("[INFO] Mitigation plan '%s' active status: %v", planID, active)
    return active, nil
}


// EnableEventMonitoring enables event monitoring in the ledger.
func (asl *AdvancedSecurityLedger) EnableEventMonitoring() error {
    asl.mu.Lock()
    defer asl.mu.Unlock()

    if asl.EventMonitoring {
        return fmt.Errorf("event monitoring is already enabled")
    }

    asl.EventMonitoring = true
    log.Printf("[INFO] Event monitoring enabled")
    return nil
}


// SetEventMonitoringStatus sets the status of event monitoring (enabled/disabled).
func (asl *AdvancedSecurityLedger) SetEventMonitoringStatus(enabled bool) error {
    asl.mu.Lock()
    defer asl.mu.Unlock()

    previousStatus := asl.EventMonitoring
    asl.EventMonitoring = enabled

    if previousStatus == enabled {
        status := "enabled"
        if !enabled {
            status = "disabled"
        }
        return fmt.Errorf("event monitoring is already %s", status)
    }

    status := "disabled"
    if enabled {
        status = "enabled"
    }

    log.Printf("[INFO] Event monitoring status changed to: %s", status)
    return nil
}


// RecordAPIUsage records API usage metrics in the ledger.
func (asl *AdvancedSecurityLedger) RecordAPIUsage(apiName string, usageCount int) error {
    if apiName == "" || usageCount < 0 {
        return fmt.Errorf("invalid API name or usage count")
    }

    asl.Lock()
    defer asl.Unlock()

    if asl.APIUsageMetrics == nil {
        asl.APIUsageMetrics = make(map[string]int)
    }

    asl.APIUsageMetrics[apiName] += usageCount
    log.Printf("[INFO] Recorded API usage: %s -> %d total requests", apiName, asl.APIUsageMetrics[apiName])
    return nil
}



// NewThreatDetectionManager creates a new ThreatDetectionManager.
func NewThreatDetectionManager() *ThreatDetectionManager {
    log.Println("[INFO] ThreatDetectionManager instance created.")
    return &ThreatDetectionManager{Active: false}
}


// DeactivateThreatDetection deactivates the threat detection system.
func (tdm *ThreatDetectionManager) DeactivateThreatDetection() error {
    if !tdm.Active {
        return fmt.Errorf("threat detection is already disabled")
    }

    tdm.Active = false
    log.Println("[INFO] Threat detection system deactivated.")
    return nil
}


// EnableIntrusionDetection activates intrusion detection systems.
func (asl *AdvancedSecurityLedger) EnableIntrusionDetection() error {
    asl.mu.Lock()
    defer asl.mu.Unlock()

    if asl.IntrusionDetectionEnabled {
        return fmt.Errorf("intrusion detection is already enabled")
    }

    asl.IntrusionDetectionEnabled = true
    log.Println("[INFO] Intrusion detection systems activated.")
    return nil
}


// DisableIntrusionDetection deactivates intrusion detection systems.
func (asl *AdvancedSecurityLedger) DisableIntrusionDetection() error {
    asl.mu.Lock()
    defer asl.mu.Unlock()

    if !asl.IntrusionDetectionEnabled {
        return fmt.Errorf("intrusion detection is already disabled")
    }

    asl.IntrusionDetectionEnabled = false
    log.Println("[INFO] Intrusion detection systems deactivated.")
    return nil
}

