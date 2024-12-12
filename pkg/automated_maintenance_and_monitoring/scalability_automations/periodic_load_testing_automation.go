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
    LoadTestInterval          = 60 * time.Minute // Interval between load tests
    MaxLoadTestRetries        = 3                // Maximum number of retries for failed load tests
    LoadTestThresholdTxPerSec = 500              // Acceptable transactions per second (TPS) threshold
)

// PeriodicLoadTestingAutomation manages the periodic load testing of the blockchain to assess performance and scalability.
type PeriodicLoadTestingAutomation struct {
    ledgerInstance    *ledger.Ledger               // Reference to the ledger for logging
    consensusSystem   *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    stateMutex        *sync.RWMutex                // Mutex for concurrency control
    loadTestAttempts  map[string]int               // Tracks retry attempts for load tests
}

// NewPeriodicLoadTestingAutomation creates and initializes a new PeriodicLoadTestingAutomation.
func NewPeriodicLoadTestingAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *PeriodicLoadTestingAutomation {
    return &PeriodicLoadTestingAutomation{
        ledgerInstance:   ledgerInstance,
        consensusSystem:  consensusSystem,
        stateMutex:       stateMutex,
        loadTestAttempts: make(map[string]int),
    }
}

// StartLoadTestingAutomation starts the load testing automation in a continuous loop.
func (automation *PeriodicLoadTestingAutomation) StartLoadTestingAutomation() {
    ticker := time.NewTicker(LoadTestInterval)
    go func() {
        for range ticker.C {
            automation.runLoadTest()
        }
    }()
}

// runLoadTest initiates a new load test to check the network's scalability and performance under load.
func (automation *PeriodicLoadTestingAutomation) runLoadTest() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    loadTestID := automation.generateLoadTestID()

    fmt.Printf("Starting load test with ID: %s\n", loadTestID)

    success, tps := automation.performLoadTest(loadTestID)
    if !success {
        automation.retryLoadTest(loadTestID)
        return
    }

    // Encrypt TPS data before logging
    encryptedTPS, err := automation.encryptTPSData(tps)
    if err != nil {
        fmt.Printf("Error encrypting TPS data: %v\n", err)
        return
    }

    // Log the load test result in the ledger
    automation.logLoadTestResult(loadTestID, "Success", encryptedTPS)

    fmt.Printf("Load test completed successfully with TPS: %d\n", tps)
}

// performLoadTest simulates a load on the network to assess performance and scalability.
func (automation *PeriodicLoadTestingAutomation) performLoadTest(loadTestID string) (bool, int) {
    fmt.Printf("Performing load test %s...\n", loadTestID)

    // Simulate transactions to assess the network's transactions per second (TPS)
    tps, err := automation.consensusSystem.SimulateLoad()
    if err != nil {
        fmt.Printf("Error simulating load test %s: %v\n", loadTestID, err)
        return false, 0
    }

    // If the TPS is below the threshold, return failure
    if tps < LoadTestThresholdTxPerSec {
        fmt.Printf("Load test %s failed. TPS below threshold: %d\n", loadTestID, tps)
        return false, tps
    }

    return true, tps
}

// retryLoadTest retries the load test if it fails, up to MaxLoadTestRetries.
func (automation *PeriodicLoadTestingAutomation) retryLoadTest(loadTestID string) {
    automation.loadTestAttempts[loadTestID]++
    if automation.loadTestAttempts[loadTestID] < MaxLoadTestRetries {
        fmt.Printf("Retrying load test %s...\n", loadTestID)
        success, tps := automation.performLoadTest(loadTestID)
        if success {
            encryptedTPS, err := automation.encryptTPSData(tps)
            if err != nil {
                fmt.Printf("Error encrypting TPS data during retry: %v\n", err)
                return
            }
            automation.logLoadTestResult(loadTestID, "Success (after retry)", encryptedTPS)
            fmt.Printf("Load test %s completed successfully on retry with TPS: %d\n", loadTestID, tps)
        } else {
            fmt.Printf("Retry failed for load test %s\n", loadTestID)
        }
    } else {
        fmt.Printf("Max retries reached for load test %s. Logging failure.\n", loadTestID)
        automation.logLoadTestResult(loadTestID, "Failed", 0)
    }
}

// logLoadTestResult logs the result of a load test into the ledger.
func (automation *PeriodicLoadTestingAutomation) logLoadTestResult(loadTestID string, status string, tps int) {
    entry := common.LedgerEntry{
        ID:        loadTestID,
        Timestamp: time.Now().Unix(),
        Type:      "Load Test",
        Status:    status,
        Details:   fmt.Sprintf("Load test %s status: %s, TPS: %d", loadTestID, status, tps),
    }

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Failed to log load test result for %s: %v\n", loadTestID, err)
    }
}

// encryptTPSData encrypts the TPS data before logging to the ledger.
func (automation *PeriodicLoadTestingAutomation) encryptTPSData(tps int) (int, error) {
    encryptedData, err := encryption.EncryptData([]byte(fmt.Sprintf("%d", tps)))
    if err != nil {
        return 0, fmt.Errorf("error encrypting TPS data: %v", err)
    }
    return int(encryptedData[0]), nil // Return only the encrypted first byte as an integer for simplicity
}

// generateLoadTestID generates a unique identifier for each load test.
func (automation *PeriodicLoadTestingAutomation) generateLoadTestID() string {
    return fmt.Sprintf("LOADTEST-%d", time.Now().UnixNano())
}
