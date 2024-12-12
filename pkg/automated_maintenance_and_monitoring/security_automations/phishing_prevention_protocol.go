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
    PhishingMonitoringInterval   = 10 * time.Second // Interval for monitoring potential phishing attacks
    MaxPhishingPreventionRetries = 3                // Maximum retries for preventing phishing attempts
    SubBlocksPerBlock            = 1000             // Number of sub-blocks in a block
    PhishingAlertThreshold       = 5                // Threshold for alerting or blocking phishing attempts
)

// PhishingPreventionProtocol manages phishing detection and prevention within the network
type PhishingPreventionProtocol struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger for logging phishing prevention-related events
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    phishingRetryCount       map[string]int               // Counter for retrying phishing prevention actions
    phishingMonitoringCycleCount int                      // Counter for phishing monitoring cycles
    phishingAttemptCounter   map[string]int               // Tracks phishing attempts by potential attackers
}

// NewPhishingPreventionProtocol initializes the automation for phishing prevention
func NewPhishingPreventionProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *PhishingPreventionProtocol {
    return &PhishingPreventionProtocol{
        consensusSystem:          consensusSystem,
        ledgerInstance:           ledgerInstance,
        stateMutex:               stateMutex,
        phishingRetryCount:       make(map[string]int),
        phishingAttemptCounter:   make(map[string]int),
        phishingMonitoringCycleCount: 0,
    }
}

// StartPhishingPreventionMonitoring starts the continuous loop for monitoring and detecting phishing attempts
func (protocol *PhishingPreventionProtocol) StartPhishingPreventionMonitoring() {
    ticker := time.NewTicker(PhishingMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorForPhishingAttempts()
        }
    }()
}

// monitorForPhishingAttempts monitors the network for signs of phishing attempts and takes action if detected
func (protocol *PhishingPreventionProtocol) monitorForPhishingAttempts() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch suspicious activities from the consensus system
    suspiciousActivities := protocol.consensusSystem.DetectSuspiciousPhishingAttempts()

    for _, activity := range suspiciousActivities {
        if protocol.isPhishingAttempt(activity) {
            fmt.Printf("Phishing attempt detected for entity %s. Taking action.\n", activity.EntityID)
            protocol.handlePhishingAttempt(activity)
        } else {
            fmt.Printf("No phishing detected for entity %s.\n", activity.EntityID)
        }
    }

    protocol.phishingMonitoringCycleCount++
    fmt.Printf("Phishing prevention cycle #%d completed.\n", protocol.phishingMonitoringCycleCount)

    if protocol.phishingMonitoringCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizePhishingPreventionCycle()
    }
}

// isPhishingAttempt checks if the detected activity matches phishing attack patterns
func (protocol *PhishingPreventionProtocol) isPhishingAttempt(activity common.SuspiciousActivity) bool {
    // Logic to determine if the activity is a phishing attempt (based on behavior, patterns, etc.)
    return activity.IsPhishingRelated
}

// handlePhishingAttempt handles a phishing attempt by taking immediate action, such as blocking the entity or issuing warnings
func (protocol *PhishingPreventionProtocol) handlePhishingAttempt(activity common.SuspiciousActivity) {
    protocol.phishingAttemptCounter[activity.EntityID]++

    if protocol.phishingAttemptCounter[activity.EntityID] >= PhishingAlertThreshold {
        fmt.Printf("Multiple phishing attempts detected for entity %s. Blocking access.\n", activity.EntityID)
        protocol.blockPhishingEntity(activity)
    } else {
        fmt.Printf("Issuing phishing warning to entity %s.\n", activity.EntityID)
        protocol.warnPhishingEntity(activity)
    }
}

// warnPhishingEntity issues a warning to an entity suspected of phishing
func (protocol *PhishingPreventionProtocol) warnPhishingEntity(activity common.SuspiciousActivity) {
    encryptedWarningData := protocol.encryptPhishingData(activity)

    // Issue a warning through the Synnergy Consensus system
    warningSuccess := protocol.consensusSystem.WarnPhishingEntity(activity.EntityID, encryptedWarningData)

    if warningSuccess {
        fmt.Printf("Phishing warning issued to entity %s.\n", activity.EntityID)
        protocol.logPhishingEvent(activity, "Warning Issued")
        protocol.resetPhishingRetry(activity.EntityID)
    } else {
        fmt.Printf("Error issuing phishing warning to entity %s. Retrying...\n", activity.EntityID)
        protocol.retryPhishingPrevention(activity)
    }
}

// blockPhishingEntity blocks access for an entity after multiple phishing attempts
func (protocol *PhishingPreventionProtocol) blockPhishingEntity(activity common.SuspiciousActivity) {
    encryptedBlockData := protocol.encryptPhishingData(activity)

    // Attempt to block the entity through the Synnergy Consensus system
    blockSuccess := protocol.consensusSystem.BlockPhishingEntity(activity.EntityID, encryptedBlockData)

    if blockSuccess {
        fmt.Printf("Phishing entity %s blocked.\n", activity.EntityID)
        protocol.logPhishingEvent(activity, "Entity Blocked")
        protocol.resetPhishingRetry(activity.EntityID)
    } else {
        fmt.Printf("Error blocking phishing entity %s. Retrying...\n", activity.EntityID)
        protocol.retryPhishingPrevention(activity)
    }
}

// retryPhishingPrevention retries the phishing prevention action if it initially fails
func (protocol *PhishingPreventionProtocol) retryPhishingPrevention(activity common.SuspiciousActivity) {
    protocol.phishingRetryCount[activity.EntityID]++
    if protocol.phishingRetryCount[activity.EntityID] < MaxPhishingPreventionRetries {
        protocol.handlePhishingAttempt(activity)
    } else {
        fmt.Printf("Max retries reached for phishing prevention action on entity %s. Action failed.\n", activity.EntityID)
        protocol.logPhishingPreventionFailure(activity)
    }
}

// resetPhishingRetry resets the retry count for phishing prevention actions on a specific entity
func (protocol *PhishingPreventionProtocol) resetPhishingRetry(entityID string) {
    protocol.phishingRetryCount[entityID] = 0
}

// finalizePhishingPreventionCycle finalizes the phishing monitoring cycle and logs the result in the ledger
func (protocol *PhishingPreventionProtocol) finalizePhishingPreventionCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizePhishingPreventionCycle()
    if success {
        fmt.Println("Phishing prevention cycle finalized successfully.")
        protocol.logPhishingPreventionCycleFinalization()
    } else {
        fmt.Println("Error finalizing phishing prevention cycle.")
    }
}

// logPhishingEvent logs a phishing-related event into the ledger
func (protocol *PhishingPreventionProtocol) logPhishingEvent(activity common.SuspiciousActivity, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("phishing-event-%s-%s", activity.EntityID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Phishing Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Entity %s triggered %s due to phishing activity.", activity.EntityID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with phishing event for entity %s.\n", activity.EntityID)
}

// logPhishingPreventionFailure logs the failure to prevent phishing into the ledger
func (protocol *PhishingPreventionProtocol) logPhishingPreventionFailure(activity common.SuspiciousActivity) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("phishing-prevention-failure-%s", activity.EntityID),
        Timestamp: time.Now().Unix(),
        Type:      "Phishing Prevention Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to prevent phishing for entity %s after maximum retries.", activity.EntityID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with phishing prevention failure for entity %s.\n", activity.EntityID)
}

// logPhishingPreventionCycleFinalization logs the finalization of a phishing prevention cycle into the ledger
func (protocol *PhishingPreventionProtocol) logPhishingPreventionCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("phishing-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Phishing Prevention Cycle Finalization",
        Status:    "Finalized",
        Details:   "Phishing prevention cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with phishing prevention cycle finalization.")
}

// encryptPhishingData encrypts the data related to phishing attempts before issuing warnings or blocking access
func (protocol *PhishingPreventionProtocol) encryptPhishingData(activity common.SuspiciousActivity) common.SuspiciousActivity {
    encryptedData, err := encryption.EncryptData(activity.ActivityData)
    if err != nil {
        fmt.Println("Error encrypting phishing data:", err)
        return activity
    }

    activity.EncryptedData = encryptedData
    fmt.Println("Phishing data successfully encrypted for entity ID:", activity.EntityID)
    return activity
}

// triggerEmergencyPhishingLockdown triggers an emergency lockdown on an entity suspected of severe phishing activity
func (protocol *PhishingPreventionProtocol) triggerEmergencyPhishingLockdown(entityID string) {
    fmt.Printf("Emergency phishing lockdown triggered for entity ID: %s.\n", entityID)
    activity := protocol.consensusSystem.GetSuspiciousActivityByID(entityID)
    encryptedData := protocol.encryptPhishingData(activity)

    success := protocol.consensusSystem.TriggerEmergencyPhishingLockdown(entityID, encryptedData)

    if success {
        protocol.logPhishingEvent(activity, "Emergency Locked Down")
        fmt.Println("Emergency phishing lockdown executed successfully.")
    } else {
        fmt.Println("Emergency phishing lockdown failed.")
    }
}
