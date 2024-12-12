package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    SmartContractEventLoggingInterval = 1500 * time.Millisecond // Interval for checking smart contract events to log
    SubBlocksPerBlock                 = 1000                    // Number of sub-blocks in a block
)

// SmartContractEventLoggingAutomation automates the logging of smart contract events
type SmartContractEventLoggingAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
    ledgerInstance        *ledger.Ledger               // Ledger to store smart contract event logs
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    loggingCycleCount     int                          // Counter for logging check cycles
}

// NewSmartContractEventLoggingAutomation initializes the automation for logging smart contract events
func NewSmartContractEventLoggingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SmartContractEventLoggingAutomation {
    return &SmartContractEventLoggingAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        loggingCycleCount: 0,
    }
}

// StartEventLogging starts the continuous loop for monitoring and logging smart contract events
func (automation *SmartContractEventLoggingAutomation) StartEventLogging() {
    ticker := time.NewTicker(SmartContractEventLoggingInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndLogEvents()
        }
    }()
}

// monitorAndLogEvents checks for smart contract events and logs them into the ledger
func (automation *SmartContractEventLoggingAutomation) monitorAndLogEvents() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the latest smart contract events from the consensus system
    events := automation.consensusSystem.GetSmartContractEvents()

    for _, event := range events {
        fmt.Printf("Logging event from contract %s.\n", event.ContractAddress)
        automation.logSmartContractEvent(event)
    }

    automation.loggingCycleCount++
    if automation.loggingCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeEventLoggingCycle()
    }
}

// logSmartContractEvent logs the specified smart contract event into the ledger for traceability
func (automation *SmartContractEventLoggingAutomation) logSmartContractEvent(event common.SmartContractEvent) {
    // Encrypt event data before logging
    encryptedEvent := automation.encryptEventData(event)

    // Log the event into the ledger
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("smart-contract-event-%s-%d", event.ContractAddress, event.EventID),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Event",
        Status:    "Logged",
        Details:   fmt.Sprintf("Event ID: %d from contract %s logged successfully.", event.EventID, event.ContractAddress),
        EncryptedData: encryptedEvent,
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with smart contract event from contract %s.\n", event.ContractAddress)
}

// finalizeEventLoggingCycle finalizes the logging cycle and logs the result in the ledger
func (automation *SmartContractEventLoggingAutomation) finalizeEventLoggingCycle() {
    success := automation.consensusSystem.FinalizeLoggingCycle()
    if success {
        fmt.Println("Smart contract event logging cycle finalized successfully.")
        automation.logLoggingCycleFinalization()
    } else {
        fmt.Println("Error finalizing smart contract event logging cycle.")
    }
}

// logLoggingCycleFinalization logs the finalization of an event logging cycle into the ledger
func (automation *SmartContractEventLoggingAutomation) logLoggingCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("event-logging-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Event Logging Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with event logging cycle finalization.")
}

// encryptEventData encrypts the event data before logging it into the ledger
func (automation *SmartContractEventLoggingAutomation) encryptEventData(event common.SmartContractEvent) []byte {
    encryptedData, err := encryption.EncryptData(event)
    if err != nil {
        fmt.Println("Error encrypting event data:", err)
        return nil
    }
    fmt.Println("Event data successfully encrypted.")
    return encryptedData
}

// ensureEventLoggingIntegrity checks the integrity of event logs and triggers revalidation if necessary
func (automation *SmartContractEventLoggingAutomation) ensureEventLoggingIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateEventLoggingIntegrity()
    if !integrityValid {
        fmt.Println("Event logging data integrity breach detected. Re-triggering event logging checks.")
        automation.monitorAndLogEvents()
    } else {
        fmt.Println("Event logging data integrity is valid.")
    }
}
