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
    OracleDataVerificationInterval  = 10 * time.Second // Interval for verifying oracle data
    MaxDataVerificationRetries      = 3                // Maximum retries for verifying data
    SubBlocksPerBlock               = 1000             // Number of sub-blocks in a block
    OracleDataConsistencyThreshold  = 90.0             // Threshold for data consistency in percentage
    OracleResponseTimeThreshold     = 5 * time.Second  // Maximum allowed response time for oracle data
)

// OracleDataVerificationProtocol ensures the accuracy and integrity of data coming from oracles
type OracleDataVerificationProtocol struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger for logging data verification-related events
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    verificationRetryCount map[string]int              // Counter for retrying data verification
    oracleVerificationCycleCount int                   // Counter for verification cycles
    verifiedOracleData    map[string]bool              // Tracks whether data from an oracle has been verified
}

// NewOracleDataVerificationProtocol initializes the automation for oracle data verification
func NewOracleDataVerificationProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *OracleDataVerificationProtocol {
    return &OracleDataVerificationProtocol{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        verificationRetryCount:  make(map[string]int),
        verifiedOracleData:      make(map[string]bool),
        oracleVerificationCycleCount: 0,
    }
}

// StartOracleDataVerificationMonitoring starts the continuous loop for monitoring oracle data verification
func (protocol *OracleDataVerificationProtocol) StartOracleDataVerificationMonitoring() {
    ticker := time.NewTicker(OracleDataVerificationInterval)

    go func() {
        for range ticker.C {
            protocol.verifyOracleData()
        }
    }()
}

// verifyOracleData checks the data from decentralized oracles for accuracy, integrity, and consistency
func (protocol *OracleDataVerificationProtocol) verifyOracleData() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch oracle data from the consensus system
    oracleDataList := protocol.consensusSystem.FetchOracleData()

    for _, oracleData := range oracleDataList {
        if protocol.isDataConsistent(oracleData) {
            fmt.Printf("Oracle data from %s is consistent and verified.\n", oracleData.OracleID)
            protocol.markDataAsVerified(oracleData)
        } else {
            fmt.Printf("Inconsistency detected in data from oracle %s. Triggering verification.\n", oracleData.OracleID)
            protocol.handleDataInconsistency(oracleData)
        }
    }

    protocol.oracleVerificationCycleCount++
    fmt.Printf("Oracle data verification cycle #%d completed.\n", protocol.oracleVerificationCycleCount)

    if protocol.oracleVerificationCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeVerificationCycle()
    }
}

// isDataConsistent checks whether the oracle data is consistent with the threshold and response time
func (protocol *OracleDataVerificationProtocol) isDataConsistent(oracleData common.OracleData) bool {
    consistencyPercentage := protocol.calculateDataConsistency(oracleData)
    responseTime := time.Since(oracleData.Timestamp)

    return consistencyPercentage >= OracleDataConsistencyThreshold && responseTime <= OracleResponseTimeThreshold
}

// calculateDataConsistency calculates the consistency percentage of the oracle data
func (protocol *OracleDataVerificationProtocol) calculateDataConsistency(oracleData common.OracleData) float64 {
    // Example logic: Calculate the consistency percentage (can be customized based on specific requirements)
    consistencyScore := protocol.consensusSystem.CalculateOracleDataConsistency(oracleData)
    return consistencyScore
}

// handleDataInconsistency handles cases where oracle data is inconsistent or not verified
func (protocol *OracleDataVerificationProtocol) handleDataInconsistency(oracleData common.OracleData) {
    protocol.verificationRetryCount[oracleData.OracleID]++

    if protocol.verificationRetryCount[oracleData.OracleID] >= MaxDataVerificationRetries {
        fmt.Printf("Oracle data from %s failed verification after maximum retries. Taking action.\n", oracleData.OracleID)
        protocol.flagDataAsUnreliable(oracleData)
    } else {
        fmt.Printf("Retrying verification for oracle data from %s.\n", oracleData.OracleID)
        protocol.retryDataVerification(oracleData)
    }
}

// retryDataVerification retries the verification of inconsistent oracle data
func (protocol *OracleDataVerificationProtocol) retryDataVerification(oracleData common.OracleData) {
    protocol.verifyOracleData() // Recheck the data after retrying logic
}

// markDataAsVerified marks the oracle data as verified and logs the result in the ledger
func (protocol *OracleDataVerificationProtocol) markDataAsVerified(oracleData common.OracleData) {
    protocol.verifiedOracleData[oracleData.OracleID] = true
    protocol.logDataVerificationEvent(oracleData, "Data Verified")
}

// flagDataAsUnreliable flags the oracle data as unreliable in the consensus system and logs the result
func (protocol *OracleDataVerificationProtocol) flagDataAsUnreliable(oracleData common.OracleData) {
    protocol.consensusSystem.MarkOracleDataAsUnreliable(oracleData.OracleID)
    protocol.logDataVerificationEvent(oracleData, "Data Unreliable")
}

// finalizeVerificationCycle finalizes the verification cycle and logs the result in the ledger
func (protocol *OracleDataVerificationProtocol) finalizeVerificationCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeDataVerificationCycle()
    if success {
        fmt.Println("Oracle data verification cycle finalized successfully.")
        protocol.logVerificationCycleFinalization()
    } else {
        fmt.Println("Error finalizing oracle data verification cycle.")
    }
}

// logDataVerificationEvent logs a data verification-related event into the ledger
func (protocol *OracleDataVerificationProtocol) logDataVerificationEvent(oracleData common.OracleData, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("oracle-data-verification-%s-%s", oracleData.OracleID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Oracle Data Verification Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Oracle data from %s was %s.", oracleData.OracleID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with oracle data verification event for oracle %s.\n", oracleData.OracleID)
}

// logVerificationCycleFinalization logs the finalization of a data verification cycle into the ledger
func (protocol *OracleDataVerificationProtocol) logVerificationCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("verification-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Verification Cycle Finalization",
        Status:    "Finalized",
        Details:   "Oracle data verification cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with oracle data verification cycle finalization.")
}

// encryptOracleData encrypts the oracle data before performing verification actions
func (protocol *OracleDataVerificationProtocol) encryptOracleData(oracleData common.OracleData) common.OracleData {
    encryptedData, err := encryption.EncryptData(oracleData.Data)
    if err != nil {
        fmt.Println("Error encrypting oracle data:", err)
        return oracleData
    }

    oracleData.EncryptedData = encryptedData
    fmt.Println("Oracle data successfully encrypted for oracle ID:", oracleData.OracleID)
    return oracleData
}

// triggerEmergencyDataLockdown triggers an emergency lockdown on oracle data in case of severe inconsistencies
func (protocol *OracleDataVerificationProtocol) triggerEmergencyDataLockdown(oracleID string) {
    fmt.Printf("Emergency data lockdown triggered for oracle ID: %s.\n", oracleID)
    oracleData := protocol.consensusSystem.GetOracleDataByID(oracleID)
    encryptedData := protocol.encryptOracleData(oracleData)

    success := protocol.consensusSystem.TriggerEmergencyDataLockdown(oracleID, encryptedData)

    if success {
        protocol.logDataVerificationEvent(oracleData, "Emergency Locked Down")
        fmt.Println("Emergency data lockdown executed successfully.")
    } else {
        fmt.Println("Emergency data lockdown failed.")
    }
}
