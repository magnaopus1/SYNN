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
    ShardActivityLoggingInterval = 20 * time.Second // Interval for logging shard activities
    ShardLogEncryptionKey        = "superSecretKey" // Encryption key for sensitive logging data
)

// ShardActivityLoggingAutomation automates the logging of shard activities in the Synnergy Network.
type ShardActivityLoggingAutomation struct {
    ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
    consensusSystem *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    shardID         string                       // Shard ID being monitored
    shardMutex      *sync.RWMutex                // Mutex for concurrency control
}

// NewShardActivityLoggingAutomation creates a new instance of ShardActivityLoggingAutomation.
func NewShardActivityLoggingAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardID string, shardMutex *sync.RWMutex) *ShardActivityLoggingAutomation {
    return &ShardActivityLoggingAutomation{
        ledgerInstance:  ledgerInstance,
        consensusSystem: consensusSystem,
        shardID:         shardID,
        shardMutex:      shardMutex,
    }
}

// StartShardActivityLogging initiates the continuous shard activity logging.
func (automation *ShardActivityLoggingAutomation) StartShardActivityLogging() {
    ticker := time.NewTicker(ShardActivityLoggingInterval)
    go func() {
        for range ticker.C {
            automation.logShardActivities()
        }
    }()
}

// logShardActivities fetches and logs activities of the shard, including validation and transaction execution.
func (automation *ShardActivityLoggingAutomation) logShardActivities() {
    automation.shardMutex.Lock()
    defer automation.shardMutex.Unlock()

    activityReport, err := automation.fetchShardActivityReport()
    if err != nil {
        fmt.Printf("Error fetching shard activity report for shard %s: %v\n", automation.shardID, err)
        return
    }

    fmt.Printf("Logging activities for shard: %s\n", automation.shardID)

    encryptedActivityDetails, err := automation.encryptActivityReport(activityReport)
    if err != nil {
        fmt.Printf("Error encrypting shard activity data: %v\n", err)
        return
    }

    automation.logActivityToLedger(activityReport, encryptedActivityDetails)
}

// fetchShardActivityReport retrieves the latest shard activity report from the consensus system.
func (automation *ShardActivityLoggingAutomation) fetchShardActivityReport() (common.ShardActivityReport, error) {
    report, err := automation.consensusSystem.GetShardActivityReport(automation.shardID)
    if err != nil {
        return common.ShardActivityReport{}, fmt.Errorf("failed to get shard activity report: %v", err)
    }
    return report, nil
}

// logActivityToLedger logs the encrypted shard activity report to the ledger.
func (automation *ShardActivityLoggingAutomation) logActivityToLedger(report common.ShardActivityReport, encryptedDetails string) {
    ledgerEntry := common.LedgerEntry{
        ID:        fmt.Sprintf("SHARD-ACTIVITY-%s-%d", automation.shardID, time.Now().UnixNano()),
        Timestamp: time.Now().Unix(),
        Type:      "Shard Activity",
        Status:    "Logged",
        Details:   fmt.Sprintf("Encrypted Shard Activity Details: %s", encryptedDetails),
    }

    if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
        fmt.Printf("Failed to log shard activity for shard %s: %v\n", automation.shardID, err)
    }
}

// encryptActivityReport encrypts the shard activity report before logging it into the ledger.
func (automation *ShardActivityLoggingAutomation) encryptActivityReport(report common.ShardActivityReport) (string, error) {
    activityData := fmt.Sprintf("Shard: %s, SubBlockCount: %d, TransactionsProcessed: %d",
        report.ShardID, report.SubBlockCount, report.TransactionsProcessed)

    encryptedData, err := encryption.EncryptDataWithKey([]byte(activityData), ShardLogEncryptionKey)
    if err != nil {
        return "", fmt.Errorf("encryption failed: %v", err)
    }

    return string(encryptedData), nil
}
