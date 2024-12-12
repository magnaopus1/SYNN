package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/sharding"
)

const (
	CrossShardCheckInterval   = 5 * time.Minute // Interval for checking cross-shard resource execution
	CrossShardResourceTimeout = 15 * time.Minute // Timeout duration for resource transfer finalization
)

// CrossShardResourceExecutionAutomation handles cross-shard resource allocation and execution.
type CrossShardResourceExecutionAutomation struct {
	consensusEngine       *synnergy_consensus.SynnergyConsensus // Consensus engine for validation across shards
	ledgerInstance        *ledger.Ledger                        // Ledger for logging cross-shard resource events
	shardManager          *sharding.ShardManager                // Manages the resources across different shards
	stateMutex            *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewCrossShardResourceExecutionAutomation initializes cross-shard resource execution automation
func NewCrossShardResourceExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, shardManager *sharding.ShardManager, stateMutex *sync.RWMutex) *CrossShardResourceExecutionAutomation {
	return &CrossShardResourceExecutionAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		shardManager:    shardManager,
		stateMutex:      stateMutex,
	}
}

// StartCrossShardExecutionMonitor begins continuous monitoring and execution of cross-shard resource transfers
func (automation *CrossShardResourceExecutionAutomation) StartCrossShardExecutionMonitor() {
	ticker := time.NewTicker(CrossShardCheckInterval)

	go func() {
		for range ticker.C {
			automation.processPendingCrossShardTransfers()
		}
	}()
}

// processPendingCrossShardTransfers checks for pending cross-shard resource transfers that need execution
func (automation *CrossShardResourceExecutionAutomation) processPendingCrossShardTransfers() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch pending cross-shard resource transfers
	pendingTransfers, err := automation.shardManager.FetchPendingCrossShardTransfers()
	if err != nil {
		fmt.Println("Error fetching pending cross-shard transfers:", err)
		return
	}

	for _, transfer := range pendingTransfers {
		// Validate the transfer across shards using Synnergy Consensus
		if automation.validateCrossShardTransfer(transfer) {
			automation.executeAndLogTransfer(transfer)
		} else {
			fmt.Printf("Cross-shard transfer %s failed consensus validation.\n", transfer.ID)
		}
	}
}

// validateCrossShardTransfer ensures the transfer passes the consensus validation across shards
func (automation *CrossShardResourceExecutionAutomation) validateCrossShardTransfer(transfer sharding.CrossShardTransfer) bool {
	// Validate transfer with Synnergy Consensus
	isValid := automation.consensusEngine.ValidateCrossShardTransfer(transfer)
	return isValid
}

// executeAndLogTransfer completes the resource transfer and logs it into the ledger
func (automation *CrossShardResourceExecutionAutomation) executeAndLogTransfer(transfer sharding.CrossShardTransfer) {
	// Execute the resource transfer between shards
	err := automation.shardManager.ExecuteCrossShardTransfer(transfer)
	if err != nil {
		fmt.Printf("Error executing cross-shard transfer %s: %v\n", transfer.ID, err)
		return
	}

	// Log the transfer execution into the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("cross-shard-transfer-%s", transfer.ID),
		Timestamp: time.Now().Unix(),
		Type:      "Cross-Shard Transfer",
		Status:    "Completed",
		Details:   fmt.Sprintf("Cross-shard resource transfer %s successfully executed.", transfer.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log cross-shard transfer %s: %v\n", transfer.ID, err)
	} else {
		fmt.Printf("Cross-shard transfer %s executed and logged successfully.\n", transfer.ID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *CrossShardResourceExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

