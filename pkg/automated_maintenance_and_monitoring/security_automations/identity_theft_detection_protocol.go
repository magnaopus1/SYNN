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
    IdentityTheftDetectionInterval = 30 * time.Second // Interval for monitoring identity theft activity
    MaxIdentityAlertRetries        = 5                // Maximum retries for alerting about identity theft
    SubBlocksPerBlock              = 1000             // Number of sub-blocks in a block
)

// IdentityTheftDetectionProtocol monitors for identity theft activities in the network
type IdentityTheftDetectionProtocol struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging identity theft-related events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    identityAlertRetries map[string]int               // Counter for retrying identity theft alerts
    detectionCycleCount  int                          // Counter for identity theft detection cycles
}

// NewIdentityTheftDetectionProtocol initializes the automation for detecting identity theft
func NewIdentityTheftDetectionProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *IdentityTheftDetectionProtocol {
    return &IdentityTheftDetectionProtocol{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        identityAlertRetries: make(map[string]int),
        detectionCycleCount:  0,
    }
}

// StartIdentityTheftMonitoring starts the continuous loop for monitoring identity theft activities
func (protocol *IdentityTheftDetectionProtocol) StartIdentityTheftMonitoring() {
    ticker := time.NewTicker(IdentityTheftDetectionInterval)

    go func() {
        for range ticker.C {
            protocol.detectAndRespondToIdentityTheft()
        }
    }()
}

// detectAndRespondToIdentityTheft checks the blockchain for identity theft indicators and responds accordingly
func (protocol *IdentityTheftDetectionProtocol) detectAndRespondToIdentityTheft() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch suspicious activities related to identity theft
    suspiciousActivities := protocol.consensusSystem.DetectIdentityTheft()

    if len(suspiciousActivities) > 0 {
        for _, activity := range suspiciousActivities {
            fmt.Printf("Suspicious identity activity detected for account %s. Investigating.\n", activity.AccountID)
            protocol.handleIdentityTheft(activity)
        }
    } else {
        fmt.Println("No suspicious identity activities detected this cycle.")
    }

    protocol.detectionCycleCount++
    fmt.Printf("Identity theft detection cycle #%d completed.\n", protocol.detectionCycleCount)

    if protocol.detectionCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeIdentityTheftCycle()
    }
}

// handleIdentityTheft handles a detected identity theft case by applying response actions
func (protocol *IdentityTheftDetectionProtocol) handleIdentityTheft(activity common.IdentityActivity) {
    encryptedIdentityData := protocol.encryptIdentityData(activity)

    // Attempt to apply protective actions through the Synnergy Consensus system
    responseSuccess := protocol.consensusSystem.RespondToIdentityTheft(activity, encryptedIdentityData)

    if responseSuccess {
        fmt.Printf("Identity theft response applied successfully for account %s.\n", activity.AccountID)
        protocol.logIdentityTheftEvent(activity, "Responded")
        protocol.resetAlertRetry(activity.AccountID)
    } else {
        fmt.Printf("Error responding to identity theft for account %s. Retrying...\n", activity.AccountID)
        protocol.retryIdentityTheftAlert(activity)
    }
}

// retryIdentityTheftAlert retries identity theft alert if initial response fails
func (protocol *IdentityTheftDetectionProtocol) retryIdentityTheftAlert(activity common.IdentityActivity) {
    protocol.identityAlertRetries[activity.AccountID]++
    if protocol.identityAlertRetries[activity.AccountID] < MaxIdentityAlertRetries {
        protocol.handleIdentityTheft(activity)
    } else {
        fmt.Printf("Max retries reached for identity theft alert on account %s. Alert failed.\n", activity.AccountID)
        protocol.logIdentityTheftFailure(activity)
    }
}

// resetAlertRetry resets the retry count for identity theft alerts
func (protocol *IdentityTheftDetectionProtocol) resetAlertRetry(accountID string) {
    protocol.identityAlertRetries[accountID] = 0
}

// finalizeIdentityTheftCycle finalizes the identity theft detection cycle and logs the result in the ledger
func (protocol *IdentityTheftDetectionProtocol) finalizeIdentityTheftCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeIdentityTheftCycle()
    if success {
        fmt.Println("Identity theft detection cycle finalized successfully.")
        protocol.logIdentityTheftCycleFinalization()
    } else {
        fmt.Println("Error finalizing identity theft detection cycle.")
    }
}

// logIdentityTheftEvent logs an identity theft event into the ledger
func (protocol *IdentityTheftDetectionProtocol) logIdentityTheftEvent(activity common.IdentityActivity, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("identity-theft-%s-%s", activity.AccountID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Identity Theft Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Account %s %s due to suspicious identity activity.", activity.AccountID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with identity theft event for account %s.\n", activity.AccountID)
}

// logIdentityTheftFailure logs the failure of an identity theft response into the ledger
func (protocol *IdentityTheftDetectionProtocol) logIdentityTheftFailure(activity common.IdentityActivity) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("identity-theft-failure-%s", activity.AccountID),
        Timestamp: time.Now().Unix(),
        Type:      "Identity Theft Response Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to respond to identity theft on account %s after maximum retries.", activity.AccountID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with identity theft response failure for account %s.\n", activity.AccountID)
}

// logIdentityTheftCycleFinalization logs the finalization of an identity theft detection cycle into the ledger
func (protocol *IdentityTheftDetectionProtocol) logIdentityTheftCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("identity-theft-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Identity Theft Cycle Finalization",
        Status:    "Finalized",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with identity theft detection cycle finalization.")
}

// encryptIdentityData encrypts identity data before applying protective actions
func (protocol *IdentityTheftDetectionProtocol) encryptIdentityData(activity common.IdentityActivity) common.IdentityActivity {
    encryptedData, err := encryption.EncryptData(activity.IdentityDetails)
    if err != nil {
        fmt.Println("Error encrypting identity data:", err)
        return activity
    }

    activity.EncryptedData = encryptedData
    fmt.Println("Identity data successfully encrypted for protective actions.")
    return activity
}

// triggerEmergencyAccountLockout triggers an emergency lockout of an account in case of critical identity theft detection
func (protocol *IdentityTheftDetectionProtocol) triggerEmergencyAccountLockout(account common.IdentityActivity) {
    fmt.Printf("Emergency lockout triggered for account %s.\n", account.AccountID)
    success := protocol.consensusSystem.TriggerEmergencyLockout(account)

    if success {
        protocol.logIdentityTheftEvent(account, "Locked Out")
        fmt.Println("Emergency lockout executed successfully.")
    } else {
        fmt.Println("Emergency lockout failed.")
    }
}
