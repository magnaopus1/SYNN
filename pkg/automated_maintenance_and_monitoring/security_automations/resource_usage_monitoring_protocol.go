package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    ResourceUsageMonitoringInterval = 10 * time.Second // Interval for monitoring resource usage
    MaxUsagePolicyRetries           = 3                // Maximum retries for enforcing resource usage policies
    SubBlocksPerBlock               = 1000             // Number of sub-blocks in a block
    UsageAnomalyThreshold           = 0.20             // Threshold for detecting anomalies in resource usage
)

// ResourceUsageMonitoringProtocol monitors and enforces resource usage policies within the network
type ResourceUsageMonitoringProtocol struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger for logging resource usage events
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    usagePolicyRetryCount    map[string]int               // Counter for retrying resource usage policy enforcement
    usageMonitoringCycleCount int                         // Counter for resource usage monitoring cycles
    usageAnomalyCounter      map[string]int               // Tracks detected usage anomalies by resource
}

// NewResourceUsageMonitoringProtocol initializes the resource usage monitoring protocol
func NewResourceUsageMonitoringProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ResourceUsageMonitoringProtocol {
    return &ResourceUsageMonitoringProtocol{
        consensusSystem:           consensusSystem,
        ledgerInstance:            ledgerInstance,
        stateMutex:                stateMutex,
        usagePolicyRetryCount:     make(map[string]int),
        usageAnomalyCounter:       make(map[string]int),
        usageMonitoringCycleCount: 0,
    }
}

// StartResourceUsageMonitoring starts the continuous loop for monitoring and managing resource usage
func (protocol *ResourceUsageMonitoringProtocol) StartResourceUsageMonitoring() {
    ticker := time.NewTicker(ResourceUsageMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorResourceUsage()
        }
    }()
}

// monitorResourceUsage checks the network for excessive or anomalous resource usage
func (protocol *ResourceUsageMonitoringProtocol) monitorResourceUsage() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch real-time resource usage data from the consensus system
    usageReports := protocol.consensusSystem.FetchResourceUsageReports()

    for _, report := range usageReports {
        if protocol.isUsageAnomalyDetected(report) {
            fmt.Printf("Resource usage anomaly detected for resource %s. Taking action.\n", report.ResourceID)
            protocol.handleUsageAnomaly(report)
        } else {
            fmt.Printf("No resource usage anomaly detected for resource %s.\n", report.ResourceID)
        }
    }

    protocol.usageMonitoringCycleCount++
    fmt.Printf("Resource usage monitoring cycle #%d completed.\n", protocol.usageMonitoringCycleCount)

    if protocol.usageMonitoringCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeUsageMonitoringCycle()
    }
}

// isUsageAnomalyDetected checks if the resource usage exceeds the threshold for anomalies
func (protocol *ResourceUsageMonitoringProtocol) isUsageAnomalyDetected(report common.ResourceUsageReport) bool {
    // Logic to detect anomalies based on usage patterns, thresholds, or deviations
    return report.UsageAnomalyScore >= UsageAnomalyThreshold
}

// handleUsageAnomaly takes action when an anomalous resource usage is detected, such as issuing warnings or restricting access
func (protocol *ResourceUsageMonitoringProtocol) handleUsageAnomaly(report common.ResourceUsageReport) {
    protocol.usageAnomalyCounter[report.ResourceID]++

    if protocol.usageAnomalyCounter[report.ResourceID] >= MaxUsagePolicyRetries {
        fmt.Printf("Multiple usage anomalies detected for resource %s. Escalating response.\n", report.ResourceID)
        protocol.escalateUsageAnomalyResponse(report)
    } else {
        fmt.Printf("Issuing alert for resource usage anomaly on resource %s.\n", report.ResourceID)
        protocol.alertForUsageAnomaly(report)
    }
}

// alertForUsageAnomaly issues an alert regarding a resource usage anomaly
func (protocol *ResourceUsageMonitoringProtocol) alertForUsageAnomaly(report common.ResourceUsageReport) {
    encryptedAlertData := protocol.encryptUsageData(report)

    // Issue an alert through the Synnergy Consensus system
    alertSuccess := protocol.consensusSystem.IssueUsageAnomalyAlert(report.ResourceID, encryptedAlertData)

    if alertSuccess {
        fmt.Printf("Resource usage anomaly alert issued for resource %s.\n", report.ResourceID)
        protocol.logUsageEvent(report, "Alert Issued")
        protocol.resetUsagePolicyRetry(report.ResourceID)
    } else {
        fmt.Printf("Error issuing resource usage anomaly alert for resource %s. Retrying...\n", report.ResourceID)
        protocol.retryUsagePolicyResponse(report)
    }
}

// escalateUsageAnomalyResponse escalates the response to an anomalous resource usage, such as limiting access or enforcing restrictions
func (protocol *ResourceUsageMonitoringProtocol) escalateUsageAnomalyResponse(report common.ResourceUsageReport) {
    encryptedEscalationData := protocol.encryptUsageData(report)

    // Attempt to limit or restrict the resource's usage through the Synnergy Consensus system
    escalationSuccess := protocol.consensusSystem.EscalateUsageAnomalyResponse(report.ResourceID, encryptedEscalationData)

    if escalationSuccess {
        fmt.Printf("Usage anomaly response escalated for resource %s.\n", report.ResourceID)
        protocol.logUsageEvent(report, "Response Escalated")
        protocol.resetUsagePolicyRetry(report.ResourceID)
    } else {
        fmt.Printf("Error escalating usage anomaly response for resource %s. Retrying...\n", report.ResourceID)
        protocol.retryUsagePolicyResponse(report)
    }
}

// retryUsagePolicyResponse retries the resource usage policy response if the initial action fails
func (protocol *ResourceUsageMonitoringProtocol) retryUsagePolicyResponse(report common.ResourceUsageReport) {
    protocol.usagePolicyRetryCount[report.ResourceID]++
    if protocol.usagePolicyRetryCount[report.ResourceID] < MaxUsagePolicyRetries {
        protocol.escalateUsageAnomalyResponse(report)
    } else {
        fmt.Printf("Max retries reached for resource usage policy response on resource %s. Response failed.\n", report.ResourceID)
        protocol.logUsageFailure(report)
    }
}

// resetUsagePolicyRetry resets the retry count for usage policy responses on a specific resource
func (protocol *ResourceUsageMonitoringProtocol) resetUsagePolicyRetry(resourceID string) {
    protocol.usagePolicyRetryCount[resourceID] = 0
}

// finalizeUsageMonitoringCycle finalizes the resource usage monitoring cycle and logs the result in the ledger
func (protocol *ResourceUsageMonitoringProtocol) finalizeUsageMonitoringCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeUsageMonitoringCycle()
    if success {
        fmt.Println("Resource usage monitoring cycle finalized successfully.")
        protocol.logUsageMonitoringCycleFinalization()
    } else {
        fmt.Println("Error finalizing resource usage monitoring cycle.")
    }
}

// logUsageEvent logs a resource usage-related event into the ledger
func (protocol *ResourceUsageMonitoringProtocol) logUsageEvent(report common.ResourceUsageReport, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("resource-usage-event-%s-%s", report.ResourceID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Resource Usage Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Resource %s triggered %s due to detected usage anomaly.", report.ResourceID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with resource usage event for resource %s.\n", report.ResourceID)
}

// logUsageFailure logs the failure to respond to a resource usage issue into the ledger
func (protocol *ResourceUsageMonitoringProtocol) logUsageFailure(report common.ResourceUsageReport) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("usage-policy-failure-%s", report.ResourceID),
        Timestamp: time.Now().Unix(),
        Type:      "Usage Policy Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to usage anomaly for resource %s after maximum retries.", report.ResourceID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with usage policy failure for resource %s.\n", report.ResourceID)
}

// logUsageMonitoringCycleFinalization logs the finalization of a resource usage monitoring cycle into the ledger
func (protocol *ResourceUsageMonitoringProtocol) logUsageMonitoringCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("usage-monitoring-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Usage Monitoring Cycle Finalization",
        Status:    "Finalized",
        Details:   "Resource usage monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with usage monitoring cycle finalization.")
}

// encryptUsageData encrypts the data related to resource usage anomalies before taking action or logging events
func (protocol *ResourceUsageMonitoringProtocol) encryptUsageData(report common.ResourceUsageReport) common.ResourceUsageReport {
    encryptedData, err := encryption.EncryptData(report.UsageData)
    if err != nil {
        fmt.Println("Error encrypting resource usage data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("Resource usage data successfully encrypted for resource ID:", report.ResourceID)
    return report
}

// triggerEmergencyUsageLockdown triggers an emergency resource usage lockdown in case of critical overuse or abuse
func (protocol *ResourceUsageMonitoringProtocol) triggerEmergencyUsageLockdown(resourceID string) {
    fmt.Printf("Emergency resource usage lockdown triggered for resource ID: %s.\n", resourceID)
    report := protocol.consensusSystem.GetResourceUsageReportByID(resourceID)
    encryptedData := protocol.encryptUsageData(report)

    success := protocol.consensusSystem.TriggerEmergencyUsageLockdown(resourceID, encryptedData)

    if success {
        protocol.logUsageEvent(report, "Emergency Locked Down")
        fmt.Println("Emergency resource usage lockdown executed successfully.")
    } else {
        fmt.Println("Emergency resource usage lockdown failed.")
    }
}
