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
    InsuranceCheckInterval     = 10 * time.Second // Interval for checking triggering insurance conditions
    MaxClaimRetries            = 3               // Maximum number of retry attempts for claim payouts
)

// DefiInsuranceActivationAutomation automates the process of monitoring and triggering DeFi insurance
type DefiInsuranceActivationAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging insurance events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    insuranceRetryMap map[string]int               // Counter for retrying failed insurance claims
}

// NewDefiInsuranceActivationAutomation initializes the automation for DeFi insurance checks
func NewDefiInsuranceActivationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DefiInsuranceActivationAutomation {
    return &DefiInsuranceActivationAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        insuranceRetryMap: make(map[string]int),
    }
}

// StartInsuranceCheck starts the continuous loop for regularly checking insurance triggering events
func (automation *DefiInsuranceActivationAutomation) StartInsuranceCheck() {
    ticker := time.NewTicker(InsuranceCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndTriggerInsurance()
        }
    }()
}

// monitorAndTriggerInsurance checks for events that trigger DeFi insurance activation
func (automation *DefiInsuranceActivationAutomation) monitorAndTriggerInsurance() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of events that may trigger insurance
    insuranceEvents := automation.consensusSystem.GetInsuranceTriggerEvents()

    if len(insuranceEvents) > 0 {
        for _, event := range insuranceEvents {
            fmt.Printf("Insurance trigger detected for event %s. Activating insurance.\n", event.ID)
            automation.triggerInsuranceClaim(event)
        }
    } else {
        fmt.Println("No insurance triggering events detected at this time.")
    }
}

// triggerInsuranceClaim activates the insurance claim based on the event and processes the payout
func (automation *DefiInsuranceActivationAutomation) triggerInsuranceClaim(event common.InsuranceEvent) {
    // Encrypt event data before processing
    encryptedEventData := automation.encryptEventData(event)

    // Trigger insurance claim through the Synnergy Consensus system
    claimSuccess := automation.consensusSystem.ProcessInsuranceClaim(event, encryptedEventData)

    if claimSuccess {
        fmt.Printf("Insurance claim for event %s processed successfully.\n", event.ID)
        automation.logInsuranceClaimEvent(event)
        automation.resetInsuranceRetry(event.ID)
    } else {
        fmt.Printf("Error processing insurance claim for event %s. Retrying...\n", event.ID)
        automation.retryInsuranceClaim(event)
    }
}

// retryInsuranceClaim attempts to retry the failed insurance claim a limited number of times
func (automation *DefiInsuranceActivationAutomation) retryInsuranceClaim(event common.InsuranceEvent) {
    automation.insuranceRetryMap[event.ID]++
    if automation.insuranceRetryMap[event.ID] < MaxClaimRetries {
        automation.triggerInsuranceClaim(event)
    } else {
        fmt.Printf("Max retries reached for insurance claim on event %s. Claim failed.\n", event.ID)
        automation.logInsuranceClaimFailure(event)
    }
}

// resetInsuranceRetry resets the retry count for an insurance claim
func (automation *DefiInsuranceActivationAutomation) resetInsuranceRetry(eventID string) {
    automation.insuranceRetryMap[eventID] = 0
}

// logInsuranceClaimEvent logs the successful insurance claim event into the ledger
func (automation *DefiInsuranceActivationAutomation) logInsuranceClaimEvent(event common.InsuranceEvent) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insurance-claim-%s", event.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Insurance Claim",
        Status:    "Processed",
        Details:   fmt.Sprintf("Insurance claim for event %s processed successfully.", event.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with insurance claim event for event %s.\n", event.ID)
}

// logInsuranceClaimFailure logs the failure of an insurance claim event into the ledger
func (automation *DefiInsuranceActivationAutomation) logInsuranceClaimFailure(event common.InsuranceEvent) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insurance-claim-failure-%s", event.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Insurance Claim Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Insurance claim failed for event %s after maximum retries.", event.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with insurance claim failure for event %s.\n", event.ID)
}

// encryptEventData encrypts the insurance event data before processing
func (automation *DefiInsuranceActivationAutomation) encryptEventData(event common.InsuranceEvent) common.InsuranceEvent {
    encryptedData, err := encryption.EncryptData(event.Data)
    if err != nil {
        fmt.Println("Error encrypting insurance event data:", err)
        return event
    }

    event.EncryptedData = encryptedData
    fmt.Println("Insurance event data successfully encrypted.")
    return event
}

// ensureInsuranceIntegrity checks the integrity of the DeFi insurance claims and data
func (automation *DefiInsuranceActivationAutomation) ensureInsuranceIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateInsuranceIntegrity()
    if !integrityValid {
        fmt.Println("Insurance data integrity breach detected. Re-triggering insurance checks.")
        automation.monitorAndTriggerInsurance()
    } else {
        fmt.Println("DeFi insurance data integrity is valid.")
    }
}
