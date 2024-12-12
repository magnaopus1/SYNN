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
    EventLoggingInterval = 3000 * time.Millisecond // Interval for logging system function events
    SubBlocksPerBlock    = 1000                    // Number of sub-blocks per block
)

// SystemFunctionEventLoggingAutomation automates the process of logging events related to system functions
type SystemFunctionEventLoggingAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store system function events
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    eventLoggingCount   int                          // Counter for event logging cycles
}

// NewSystemFunctionEventLoggingAutomation initializes the automation for system function event logging
func NewSystemFunctionEventLoggingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SystemFunctionEventLoggingAutomation {
    return &SystemFunctionEventLoggingAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        eventLoggingCount: 0,
    }
}

// StartEventLogging starts the continuous loop for monitoring and logging system function events
func (automation *SystemFunctionEventLoggingAutomation) StartEventLogging() {
    ticker := time.NewTicker(EventLoggingInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndLogEvents()
        }
    }()
}

// monitorAndLogEvents checks for system function events and logs them into the ledger
func (automation *SystemFunctionEventLoggingAutomation) monitorAndLogEvents() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of recent system function events
    events, err := automation.consensusSystem.GetRecentFunctionEvents()
    if err != nil {
        fmt.Printf("Error fetching system function events: %v\n", err)
        return
    }

    // Process each event and log it in the ledger
    for _, event := range events {
        fmt.Printf("Logging event for function: %s\n", event.FunctionID)

        // Encrypt event data before logging
        encryptedEvent, err := automation.encryptEventData(event)
        if err != nil {
            fmt.Printf("Error encrypting data for function %s: %v\n", event.FunctionID, err)
            automation.logEventToLedger(event, "Encryption Failed")
            continue
        }

        automation.logEventToLedger(encryptedEvent, "Logged Successfully")
    }

    automation.eventLoggingCount++
    fmt.Printf("System function event logging cycle #%d completed.\n", automation.eventLoggingCount)

    if automation.eventLoggingCount%SubBlocksPerBlock == 0 {
        automation.finalizeEventLoggingCycle()
    }
}

// encryptEventData encrypts the event data before logging it to the ledger
func (automation *SystemFunctionEventLoggingAutomation) encryptEventData(event common.FunctionEvent) (common.FunctionEvent, error) {
    fmt.Println("Encrypting function event data.")

    encryptedData, err := encryption.EncryptData(event)
    if err != nil {
        return event, fmt.Errorf("failed to encrypt event data: %v", err)
    }

    event.EncryptedData = encryptedData
    fmt.Println("Function event data successfully encrypted.")
    return event, nil
}

// logEventToLedger logs the event in the ledger for auditability
func (automation *SystemFunctionEventLoggingAutomation) logEventToLedger(event common.FunctionEvent, status string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("function-event-%s", event.FunctionID),
        Timestamp: time.Now().Unix(),
        Type:      "System Function Event",
        Status:    status,
        Details:   fmt.Sprintf("Event for function %s: %s", event.FunctionID, status),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with event for function %s: %s\n", event.FunctionID, status)
}

// finalizeEventLoggingCycle finalizes the event logging cycle and logs the result in the ledger
func (automation *SystemFunctionEventLoggingAutomation) finalizeEventLoggingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeEventLoggingCycle()
    if success {
        fmt.Println("System function event logging cycle finalized successfully.")
        automation.logEventLoggingCycleFinalization()
    } else {
        fmt.Println("Error finalizing system function event logging cycle.")
    }
}

// logEventLoggingCycleFinalization logs the finalization of the event logging cycle in the ledger
func (automation *SystemFunctionEventLoggingAutomation) logEventLoggingCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("event-logging-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Event Logging Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with system function event logging cycle finalization.")
}
