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
	ShardAdjustmentInterval        = 30 * time.Second  // Interval for checking shard adjustment
	ShardAdjustmentLogEncryptionKey = "adjustmentKey"  // Encryption key for adjustment logs
	LoadThreshold                   = 80               // Threshold for shard load balancing (%)
	SubBlockThreshold               = 1000             // Threshold for maximum sub-blocks before adjustment
)

// ShardAdjustmentAutomation automates shard rebalancing and adjustment based on performance.
type ShardAdjustmentAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
	consensusSystem *consensus.SynnergyConsensus // Reference to Synnergy Consensus
	shardMutex      *sync.RWMutex                // Mutex for concurrency control
}

// NewShardAdjustmentAutomation creates a new instance of ShardAdjustmentAutomation.
func NewShardAdjustmentAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardAdjustmentAutomation {
	return &ShardAdjustmentAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardAdjustment initiates the continuous monitoring and adjustment of shards.
func (automation *ShardAdjustmentAutomation) StartShardAdjustment() {
	ticker := time.NewTicker(ShardAdjustmentInterval)
	go func() {
		for range ticker.C {
			automation.adjustShards()
		}
	}()
}

// adjustShards monitors shard activity and makes adjustments to ensure balance.
func (automation *ShardAdjustmentAutomation) adjustShards() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	shardReports, err := automation.consensusSystem.GetAllShardActivityReports()
	if err != nil {
		fmt.Printf("Error fetching shard activity reports: %v\n", err)
		return
	}

	for _, report := range shardReports {
		if automation.shouldAdjustShard(report) {
			automation.performShardAdjustment(report)
		}
	}
}

// shouldAdjustShard checks if a shard needs adjustment based on load and activity.
func (automation *ShardAdjustmentAutomation) shouldAdjustShard(report common.ShardActivityReport) bool {
	isOverloaded := report.LoadPercentage >= LoadThreshold
	hasExceededSubBlocks := report.SubBlockCount >= SubBlockThreshold
	return isOverloaded || hasExceededSubBlocks
}

// performShardAdjustment rebalances the load or redistributes sub-blocks in an overloaded shard.
func (automation *ShardAdjustmentAutomation) performShardAdjustment(report common.ShardActivityReport) {
	// Adjust shard by redistributing load or merging/splitting sub-blocks
	fmt.Printf("Performing shard adjustment for shard %s\n", report.ShardID)

	err := automation.consensusSystem.AdjustShardLoad(report.ShardID)
	if err != nil {
		fmt.Printf("Error adjusting shard %s: %v\n", report.ShardID, err)
		return
	}

	// Log adjustment to the ledger
	encryptedLog, err := automation.encryptAdjustmentLog(report)
	if err != nil {
		fmt.Printf("Error encrypting shard adjustment log: %v\n", err)
		return
	}

	automation.logAdjustmentToLedger(report, encryptedLog)
}

// logAdjustmentToLedger logs the shard adjustment activity to the ledger.
func (automation *ShardAdjustmentAutomation) logAdjustmentToLedger(report common.ShardActivityReport, encryptedLog string) {
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("SHARD-ADJUST-%s-%d", report.ShardID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Adjustment",
		Status:    "Completed",
		Details:   fmt.Sprintf("Encrypted Adjustment Log: %s", encryptedLog),
	}

	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log shard adjustment for shard %s: %v\n", report.ShardID, err)
	}
}

// encryptAdjustmentLog encrypts the adjustment details for secure ledger entry.
func (automation *ShardAdjustmentAutomation) encryptAdjustmentLog(report common.ShardActivityReport) (string, error) {
	adjustmentData := fmt.Sprintf("Shard: %s, Load: %d%%, SubBlockCount: %d",
		report.ShardID, report.LoadPercentage, report.SubBlockCount)

	encryptedData, err := encryption.EncryptDataWithKey([]byte(adjustmentData), ShardAdjustmentLogEncryptionKey)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}

	return string(encryptedData), nil
}
