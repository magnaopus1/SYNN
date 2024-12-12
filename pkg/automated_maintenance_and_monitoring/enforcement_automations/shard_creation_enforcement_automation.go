package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
)

// Configuration for shard creation enforcement automation
const (
	ShardCheckInterval             = 30 * time.Second // Interval to check shard requirements
	MaxTransactionsPerShard        = 100000           // Max allowed transactions per shard before splitting
	MinTransactionsPerShard        = 30000            // Min transactions to justify shard existence
	ShardCreationViolationThreshold = 3               // Allowed violations before shard adjustment
)

// ShardCreationEnforcementAutomation monitors and enforces shard creation and optimization
type ShardCreationEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	shardTransactionMap  map[string]int // Tracks transaction volume for each shard
	shardViolationMap    map[string]int // Tracks shard threshold violations for adjustment
}

// NewShardCreationEnforcementAutomation initializes the shard creation enforcement automation
func NewShardCreationEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ShardCreationEnforcementAutomation {
	return &ShardCreationEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		shardTransactionMap:  make(map[string]int),
		shardViolationMap:    make(map[string]int),
	}
}

// StartShardCreationEnforcement begins continuous monitoring and enforcement of shard creation and adjustment
func (automation *ShardCreationEnforcementAutomation) StartShardCreationEnforcement() {
	ticker := time.NewTicker(ShardCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkShardRequirements()
		}
	}()
}

// checkShardRequirements monitors each shard’s transaction volume and enforces creation or merging as needed
func (automation *ShardCreationEnforcementAutomation) checkShardRequirements() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateShardLoad()
	automation.adjustShards()
}

// evaluateShardLoad checks each shard’s transaction volume and flags shards for adjustment
func (automation *ShardCreationEnforcementAutomation) evaluateShardLoad() {
	for _, shardID := range automation.networkManager.GetAllShards() {
		transactionCount := automation.networkManager.GetShardTransactionVolume(shardID)
		automation.shardTransactionMap[shardID] = transactionCount

		if transactionCount > MaxTransactionsPerShard {
			automation.shardViolationMap[shardID]++
		} else if transactionCount < MinTransactionsPerShard {
			automation.shardViolationMap[shardID]++
		} else {
			automation.shardViolationMap[shardID] = 0 // Reset if within limits
		}
	}
}

// adjustShards creates or merges shards based on threshold violations
func (automation *ShardCreationEnforcementAutomation) adjustShards() {
	for shardID, violations := range automation.shardViolationMap {
		transactionCount := automation.shardTransactionMap[shardID]

		if violations >= ShardCreationViolationThreshold {
			if transactionCount > MaxTransactionsPerShard {
				fmt.Printf("Shard creation enforcement triggered for shard %s due to excessive transaction load.\n", shardID)
				automation.createNewShard(shardID)
			} else if transactionCount < MinTransactionsPerShard {
				fmt.Printf("Shard merging enforcement triggered for shard %s due to low transaction volume.\n", shardID)
				automation.mergeShard(shardID)
			}
		}
	}
}

// createNewShard splits an overloaded shard to maintain performance and balance transaction volume
func (automation *ShardCreationEnforcementAutomation) createNewShard(overloadedShardID string) {
	newShardID, err := automation.networkManager.CreateShard(overloadedShardID)
	if err != nil {
		fmt.Printf("Failed to create a new shard from shard %s: %v\n", overloadedShardID, err)
		automation.logShardAction(overloadedShardID, "Shard Creation Failed", fmt.Sprintf("Original Load: %d", automation.shardTransactionMap[overloadedShardID]))
	} else {
		fmt.Printf("New shard %s created from shard %s.\n", newShardID, overloadedShardID)
		automation.logShardAction(overloadedShardID, "New Shard Created", fmt.Sprintf("New Shard ID: %s, Original Load: %d", newShardID, automation.shardTransactionMap[overloadedShardID]))
	}
}

// mergeShard consolidates a low-volume shard to optimize network resources
func (automation *ShardCreationEnforcementAutomation) mergeShard(lowVolumeShardID string) {
	err := automation.networkManager.MergeShard(lowVolumeShardID)
	if err != nil {
		fmt.Printf("Failed to merge shard %s: %v\n", lowVolumeShardID, err)
		automation.logShardAction(lowVolumeShardID, "Shard Merge Failed", fmt.Sprintf("Transaction Count: %d", automation.shardTransactionMap[lowVolumeShardID]))
	} else {
		fmt.Printf("Shard %s successfully merged to optimize resources.\n", lowVolumeShardID)
		automation.logShardAction(lowVolumeShardID, "Shard Merged", fmt.Sprintf("Transaction Count: %d", automation.shardTransactionMap[lowVolumeShardID]))
	}
}

// logShardAction securely logs actions related to shard creation and management
func (automation *ShardCreationEnforcementAutomation) logShardAction(shardID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Shard ID: %s, Details: %s", action, shardID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("shard-management-enforcement-%s-%d", shardID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Management Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log shard management enforcement action for shard %s: %v\n", shardID, err)
	} else {
		fmt.Println("Shard management enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ShardCreationEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualShardAdjustment allows administrators to manually create or merge shards based on load requirements
func (automation *ShardCreationEnforcementAutomation) TriggerManualShardAdjustment() {
	fmt.Println("Manually triggering shard adjustment for load optimization.")

	for shardID, transactionCount := range automation.shardTransactionMap {
		if transactionCount > MaxTransactionsPerShard {
			automation.createNewShard(shardID)
		} else if transactionCount < MinTransactionsPerShard {
			automation.mergeShard(shardID)
		}
	}
	automation.logShardAction("Network", "Manual Shard Adjustment Triggered", "Administrator Action")
}
