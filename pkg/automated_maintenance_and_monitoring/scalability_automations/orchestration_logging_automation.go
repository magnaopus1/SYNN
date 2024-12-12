package automations

import (
    "fmt"
    "time"
    "sync"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    LoggingInterval          = 15 * time.Second   // Interval for logging orchestration events
    MaxLoggingRetries        = 3                  // Maximum retry attempts for failed logging
    LoggingTimeout           = 10 * time.Second   // Timeout for logging operations
)

// OrchestrationLoggingAutomation manages orchestration event logging.
type OrchestrationLoggingAutomation struct {
    ledgerInstance        *ledger.Ledger               // Reference to the ledger for storing logs
    consensusSystem       *consensus.SynnergyConsensus // Reference to consensus for retrieving orchestration data
    stateMutex            *sync.RWMutex                // Mutex for thread-safe logging
    loggingAttempts       map[string]int               // Tracks logging retries for each orchestration event
}

// NewOrchestrationLoggingAutomation initializes the orchestration logging automation.
func NewOrchestrationLoggingAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *OrchestrationLoggingAutomation {
    return &OrchestrationLoggingAutomation{
        ledgerInstance:  ledgerInstance,
        consensusSystem: consensusSystem,
        stateMutex:      stateMutex,
        loggingAttempts: make(map[string]int),
    }
}

// StartLogging begins the orchestration logging automation in a continuous loop.
func (automation *OrchestrationLoggingAutomation) StartLogging() {
    ticker := time.NewTicker(LoggingInterval)
    go func() {
        for range ticker.C {
            automation.logOrchestrationEvents()
        }
    }()
}

// logOrchestrationEvents retrieves orchestration events and logs them into the ledger.
func (automation *OrchestrationLoggingAutomation) logOrchestrationEvents() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    orchestrationEvents := automation.consensusSystem.GetOrchestrationEvents()

    for _, event := range orchestrationEvents {
        if err := automation.logEvent(event); err != nil {
            fmt.Printf("Failed to log orchestration event %s. Retrying...\n", event.ID)
            automation.retryLogging(event)
        }
    }
}

// logEvent logs a single orchestration event into the ledger, encrypting the data before storing.
func (automation *OrchestrationLoggingAutomation) logEvent(event common.OrchestrationEvent) error {
    encryptedDetails, err := automation.encryptOrchestrationDetails(event.Details)
    if err != nil {
        fmt.Printf("Error encrypting event details for event %s: %v\n", event.ID, err)
        return err
    }

    entry := common.LedgerEntry{
        ID:        event.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Orchestration",
        Status:    event.Status,
        Details:   encryptedDetails,
    }

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Failed to add orchestration event %s to ledger: %v\n", event.ID, err)
        return err
    }

    fmt.Printf("Orchestration event %s successfully logged.\n", event.ID)
    return nil
}

// retryLogging handles retrying the logging of an orchestration event in case of failure.
func (automation *OrchestrationLoggingAutomation) retryLogging(event common.OrchestrationEvent) {
    automation.loggingAttempts[event.ID]++

    if automation.loggingAttempts[event.ID] < MaxLoggingRetries {
        if err := automation.logEvent(event); err != nil {
            fmt.Printf("Retry failed for logging orchestration event %s.\n", event.ID)
        }
    } else {
        fmt.Printf("Max retries reached for orchestration event %s. Marking as failed.\n", event.ID)
        automation.logFailedEvent(event)
    }
}

// logFailedEvent logs a failed orchestration event after exceeding retry attempts.
func (automation *OrchestrationLoggingAutomation) logFailedEvent(event common.OrchestrationEvent) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("failed-%s", event.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Orchestration Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to log orchestration event %s after retries.", event.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Failed orchestration event %s recorded.\n", event.ID)
}

// encryptOrchestrationDetails encrypts the orchestration details before storing in the ledger.
func (automation *OrchestrationLoggingAutomation) encryptOrchestrationDetails(details string) (string, error) {
    encryptedDetails, err := encryption.EncryptData([]byte(details))
    if err != nil {
        return "", fmt.Errorf("error encrypting orchestration details: %v", err)
    }
    return string(encryptedDetails), nil
}

// logEmergencyOrchestration handles logging emergency orchestration events separately.
func (automation *OrchestrationLoggingAutomation) logEmergencyOrchestration(shard string) {
    fmt.Printf("Emergency orchestration event triggered for shard %s.\n", shard)

    event := common.OrchestrationEvent{
        ID:        fmt.Sprintf("emergency-orchestration-%s", shard),
        Timestamp: time.Now().Unix(),
        Status:    "Emergency Triggered",
        Details:   fmt.Sprintf("Emergency orchestration triggered for shard %s", shard),
    }

    if err := automation.logEvent(event); err != nil {
        fmt.Printf("Failed to log emergency orchestration for shard %s: %v\n", shard, err)
        automation.retryLogging(event)
    }
}
