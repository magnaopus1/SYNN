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
    OracleDataCheckInterval     = 15 * time.Second // Interval for verifying oracle data
    MaxOracleRetryAttempts      = 3               // Maximum retries for oracle data verification
)

// DefiOracleDataVerificationAutomation automates the process of verifying DeFi oracle data
type DefiOracleDataVerificationAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging oracle verification events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    oracleRetryMap    map[string]int               // Counter for retrying failed oracle data verification
}

// NewDefiOracleDataVerificationAutomation initializes the automation for DeFi oracle data verification
func NewDefiOracleDataVerificationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DefiOracleDataVerificationAutomation {
    return &DefiOracleDataVerificationAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        oracleRetryMap:    make(map[string]int),
    }
}

// StartOracleDataVerification starts the continuous loop for verifying oracle data at intervals
func (automation *DefiOracleDataVerificationAutomation) StartOracleDataVerification() {
    ticker := time.NewTicker(OracleDataCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndVerifyOracleData()
        }
    }()
}

// monitorAndVerifyOracleData checks and verifies data from oracles
func (automation *DefiOracleDataVerificationAutomation) monitorAndVerifyOracleData() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch oracle data that needs verification
    oracleDataList := automation.consensusSystem.GetOracleDataForVerification()

    if len(oracleDataList) > 0 {
        for _, oracleData := range oracleDataList {
            fmt.Printf("Verifying oracle data from source %s.\n", oracleData.Source)
            automation.verifyOracleData(oracleData)
        }
    } else {
        fmt.Println("No oracle data to verify at this time.")
    }
}

// verifyOracleData performs data verification for a given oracle data point
func (automation *DefiOracleDataVerificationAutomation) verifyOracleData(oracleData common.OracleData) {
    // Encrypt oracle data before verification
    encryptedOracleData := automation.encryptOracleData(oracleData)

    // Trigger data verification through the Synnergy Consensus system
    verificationSuccess := automation.consensusSystem.VerifyOracleData(oracleData, encryptedOracleData)

    if verificationSuccess {
        fmt.Printf("Oracle data from source %s verified successfully.\n", oracleData.Source)
        automation.logOracleVerificationEvent(oracleData)
        automation.resetOracleRetry(oracleData.ID)
    } else {
        fmt.Printf("Error verifying oracle data from source %s. Retrying...\n", oracleData.Source)
        automation.retryOracleVerification(oracleData)
    }
}

// retryOracleVerification retries failed oracle data verification attempts
func (automation *DefiOracleDataVerificationAutomation) retryOracleVerification(oracleData common.OracleData) {
    automation.oracleRetryMap[oracleData.ID]++
    if automation.oracleRetryMap[oracleData.ID] < MaxOracleRetryAttempts {
        automation.verifyOracleData(oracleData)
    } else {
        fmt.Printf("Max retries reached for oracle data verification from source %s. Verification failed.\n", oracleData.Source)
        automation.logOracleVerificationFailure(oracleData)
    }
}

// resetOracleRetry resets the retry count for oracle data verification
func (automation *DefiOracleDataVerificationAutomation) resetOracleRetry(oracleDataID string) {
    automation.oracleRetryMap[oracleDataID] = 0
}

// logOracleVerificationEvent logs the successful verification of oracle data into the ledger
func (automation *DefiOracleDataVerificationAutomation) logOracleVerificationEvent(oracleData common.OracleData) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("oracle-verification-%s", oracleData.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Oracle Data Verification",
        Status:    "Verified",
        Details:   fmt.Sprintf("Oracle data from source %s successfully verified.", oracleData.Source),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with oracle verification event for source %s.\n", oracleData.Source)
}

// logOracleVerificationFailure logs the failure of oracle data verification into the ledger
func (automation *DefiOracleDataVerificationAutomation) logOracleVerificationFailure(oracleData common.OracleData) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("oracle-verification-failure-%s", oracleData.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Oracle Data Verification Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Verification failed for oracle data from source %s after maximum retries.", oracleData.Source),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with oracle verification failure for source %s.\n", oracleData.Source)
}

// encryptOracleData encrypts the oracle data before verification
func (automation *DefiOracleDataVerificationAutomation) encryptOracleData(oracleData common.OracleData) common.OracleData {
    encryptedData, err := encryption.EncryptData(oracleData.Data)
    if err != nil {
        fmt.Println("Error encrypting oracle data:", err)
        return oracleData
    }

    oracleData.EncryptedData = encryptedData
    fmt.Println("Oracle data successfully encrypted.")
    return oracleData
}

// ensureOracleDataIntegrity checks the integrity of oracle data being processed
func (automation *DefiOracleDataVerificationAutomation) ensureOracleDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateOracleDataIntegrity()
    if !integrityValid {
        fmt.Println("Oracle data integrity breach detected. Re-triggering data verification.")
        automation.monitorAndVerifyOracleData()
    } else {
        fmt.Println("Oracle data integrity is valid.")
    }
}
