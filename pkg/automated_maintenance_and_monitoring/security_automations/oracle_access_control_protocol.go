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
    OracleAccessMonitoringInterval   = 10 * time.Second // Interval for monitoring oracle access
    MaxAccessControlRetries          = 3                // Maximum retries for enforcing access control
    SubBlocksPerBlock                = 1000             // Number of sub-blocks in a block
    UnauthorizedAccessAlertThreshold = 3                // Threshold for alerting on unauthorized oracle access attempts
    OracleKeyExpirationPeriod        = 30 * 24 * time.Hour // Period for oracle access keys expiration
)

// OracleAccessControlProtocol manages access control to decentralized oracles
type OracleAccessControlProtocol struct {
    consensusSystem         *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance          *ledger.Ledger               // Ledger for logging oracle access-related events
    stateMutex              *sync.RWMutex                // Mutex for thread-safe access
    accessRetryCount        map[string]int               // Counter for retrying access control enforcement
    oracleAccessCycleCount  int                          // Counter for monitoring cycles
    oracleAccessKeys        map[string]time.Time         // Tracks oracle access keys and their expiration times
    unauthorizedAccessCount map[string]int               // Tracks unauthorized access attempts per entity
}

// NewOracleAccessControlProtocol initializes the automation for oracle access control
func NewOracleAccessControlProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *OracleAccessControlProtocol {
    return &OracleAccessControlProtocol{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        accessRetryCount:       make(map[string]int),
        oracleAccessKeys:       make(map[string]time.Time),
        unauthorizedAccessCount: make(map[string]int),
        oracleAccessCycleCount: 0,
    }
}

// StartOracleAccessMonitoring starts the continuous loop for monitoring oracle access
func (protocol *OracleAccessControlProtocol) StartOracleAccessMonitoring() {
    ticker := time.NewTicker(OracleAccessMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorOracleAccess()
        }
    }()
}

// monitorOracleAccess monitors the access to decentralized oracles and enforces access control policies
func (protocol *OracleAccessControlProtocol) monitorOracleAccess() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch current oracle access activity from the consensus system
    oracleAccessActivities := protocol.consensusSystem.FetchOracleAccessActivities()

    for _, activity := range oracleAccessActivities {
        if protocol.isAccessUnauthorized(activity) {
            fmt.Printf("Unauthorized access detected for entity %s. Triggering alert.\n", activity.EntityID)
            protocol.handleUnauthorizedAccess(activity)
        } else {
            fmt.Printf("Authorized access for entity %s to oracle %s.\n", activity.EntityID, activity.OracleID)
            protocol.logOracleAccess(activity, "Authorized Access")
        }
    }

    protocol.oracleAccessCycleCount++
    fmt.Printf("Oracle access control cycle #%d completed.\n", protocol.oracleAccessCycleCount)

    if protocol.oracleAccessCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeOracleAccessCycle()
    }
}

// isAccessUnauthorized checks if the access attempt to the oracle is unauthorized
func (protocol *OracleAccessControlProtocol) isAccessUnauthorized(activity common.OracleAccessActivity) bool {
    accessKeyExpiration, exists := protocol.oracleAccessKeys[activity.OracleID]
    if !exists || time.Now().After(accessKeyExpiration) {
        return true
    }
    return false
}

// handleUnauthorizedAccess handles unauthorized access attempts to oracles
func (protocol *OracleAccessControlProtocol) handleUnauthorizedAccess(activity common.OracleAccessActivity) {
    protocol.unauthorizedAccessCount[activity.EntityID]++

    if protocol.unauthorizedAccessCount[activity.EntityID] >= UnauthorizedAccessAlertThreshold {
        fmt.Printf("Multiple unauthorized access attempts by entity %s. Taking action.\n", activity.EntityID)
        protocol.blockEntityAccess(activity)
    } else {
        fmt.Printf("Unauthorized access warning for entity %s.\n", activity.EntityID)
        protocol.logOracleAccess(activity, "Unauthorized Access")
    }
}

// blockEntityAccess blocks access to an oracle for an entity after multiple unauthorized attempts
func (protocol *OracleAccessControlProtocol) blockEntityAccess(activity common.OracleAccessActivity) {
    encryptedAccessData := protocol.encryptAccessData(activity)

    // Attempt to block the entity's access through the Synnergy Consensus system
    blockSuccess := protocol.consensusSystem.BlockEntityFromOracle(activity.EntityID, encryptedAccessData)

    if blockSuccess {
        fmt.Printf("Access blocked for entity %s.\n", activity.EntityID)
        protocol.logOracleAccess(activity, "Access Blocked")
        protocol.resetAccessRetry(activity.EntityID)
    } else {
        fmt.Printf("Error blocking access for entity %s. Retrying...\n", activity.EntityID)
        protocol.retryAccessControl(activity)
    }
}

// retryAccessControl retries the enforcement of access control in case of failure
func (protocol *OracleAccessControlProtocol) retryAccessControl(activity common.OracleAccessActivity) {
    protocol.accessRetryCount[activity.EntityID]++
    if protocol.accessRetryCount[activity.EntityID] < MaxAccessControlRetries {
        protocol.blockEntityAccess(activity)
    } else {
        fmt.Printf("Max retries reached for blocking access for entity %s. Action failed.\n", activity.EntityID)
        protocol.logAccessControlFailure(activity)
    }
}

// resetAccessRetry resets the retry count for an entity's access control action
func (protocol *OracleAccessControlProtocol) resetAccessRetry(entityID string) {
    protocol.accessRetryCount[entityID] = 0
}

// finalizeOracleAccessCycle finalizes the oracle access monitoring cycle and logs the result in the ledger
func (protocol *OracleAccessControlProtocol) finalizeOracleAccessCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeOracleAccessCycle()
    if success {
        fmt.Println("Oracle access control cycle finalized successfully.")
        protocol.logOracleAccessCycleFinalization()
    } else {
        fmt.Println("Error finalizing oracle access control cycle.")
    }
}

// logOracleAccess logs an oracle access-related event into the ledger
func (protocol *OracleAccessControlProtocol) logOracleAccess(activity common.OracleAccessActivity, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("oracle-access-%s-%s", activity.EntityID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Oracle Access Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Entity %s had %s to oracle %s.", activity.EntityID, eventType, activity.OracleID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with oracle access event for entity %s.\n", activity.EntityID)
}

// logAccessControlFailure logs the failure of an access control action into the ledger
func (protocol *OracleAccessControlProtocol) logAccessControlFailure(activity common.OracleAccessActivity) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("oracle-access-failure-%s", activity.EntityID),
        Timestamp: time.Now().Unix(),
        Type:      "Oracle Access Control Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to block access for entity %s after maximum retries.", activity.EntityID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with oracle access control failure for entity %s.\n", activity.EntityID)
}

// logOracleAccessCycleFinalization logs the finalization of an oracle access control cycle into the ledger
func (protocol *OracleAccessControlProtocol) logOracleAccessCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("oracle-access-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Oracle Access Cycle Finalization",
        Status:    "Finalized",
        Details:   "Oracle access control cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with oracle access control cycle finalization.")
}

// encryptAccessData encrypts the oracle access data before blocking or responding to unauthorized access attempts
func (protocol *OracleAccessControlProtocol) encryptAccessData(activity common.OracleAccessActivity) common.OracleAccessActivity {
    encryptedData, err := encryption.EncryptData(activity.AccessData)
    if err != nil {
        fmt.Println("Error encrypting oracle access data:", err)
        return activity
    }

    activity.EncryptedAccessData = encryptedData
    fmt.Println("Oracle access data successfully encrypted for entity ID:", activity.EntityID)
    return activity
}

// triggerEmergencyOracleAccessLockdown triggers an emergency lockdown of oracle access for critical breaches
func (protocol *OracleAccessControlProtocol) triggerEmergencyOracleAccessLockdown(entityID string) {
    fmt.Printf("Emergency lockdown triggered for entity ID: %s.\n", entityID)
    activity := protocol.consensusSystem.GetOracleAccessActivityByID(entityID)
    encryptedData := protocol.encryptAccessData(activity)

    success := protocol.consensusSystem.TriggerEmergencyOracleAccessLockdown(entityID, encryptedData)

    if success {
        protocol.logOracleAccess(activity, "Emergency Locked Down")
        fmt.Println("Emergency oracle access lockdown executed successfully.")
    } else {
        fmt.Println("Emergency oracle access lockdown failed.")
    }
}
