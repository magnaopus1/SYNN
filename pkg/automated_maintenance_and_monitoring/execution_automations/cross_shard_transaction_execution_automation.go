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
	CrossShardTransactionCheckInterval = 5 * time.Minute // Interval for checking cross-shard transactions
	CrossShardTransactionTimeout       = 15 * time.Minute // Timeout for transaction finalization
)

// CrossShardTransactionExecutionAutomation manages cross-shard transaction execution
type CrossShardTransactionExecutionAutomation struct {
	consensusEngine    *synnergy_consensus.SynnergyConsensus // Synnergy Consensus for transaction validation
	ledgerInstance     *ledger.Ledger                        // Ledger for logging cross-shard transactions
	shardManager       *sharding.ShardManager                // Manages shard-level transactions
	transactionMutex   *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewCrossShardTransactionExecutionAutomation initializes cross-shard transaction automation
func NewCrossShardTransactionExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, shardManager *sharding.ShardManager, transactionMutex *sync.RWMutex) *CrossShardTransactionExecutionAutomation {
	return &CrossShardTransactionExecutionAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		shardManager:     shardManager,
		transactionMutex: transactionMutex,
	}
}

// StartCrossShardTransactionMonitor begins the process of monitoring and executing cross-shard transactions
func (automation *CrossShardTransactionExecutionAutomation) StartCrossShardTransactionMonitor() {
	ticker := time.NewTicker(CrossShardTransactionCheckInterval)

	go func() {
		for range ticker.C {
			automation.processPendingCrossShardTransactions()
		}
	}()
}

// processPendingCrossShardTransactions retrieves and processes any pending cross-shard transactions
func (automation *CrossShardTransactionExecutionAutomation) processPendingCrossShardTransactions() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	// Fetch pending cross-shard transactions
	pendingTransactions, err := automation.shardManager.FetchPendingCrossShardTransactions()
	if err != nil {
		fmt.Println("Error fetching pending cross-shard transactions:", err)
		return
	}

	for _, transaction := range pendingTransactions {
		// Validate the transaction with Synnergy Consensus
		if automation.validateCrossShardTransaction(transaction) {
			automation.executeAndLogTransaction(transaction)
		} else {
			fmt.Printf("Cross-shard transaction %s failed consensus validation.\n", transaction.ID)
		}
	}
}

// validateCrossShardTransaction validates the cross-shard transaction using Synnergy Consensus
func (automation *CrossShardTransactionExecutionAutomation) validateCrossShardTransaction(transaction sharding.CrossShardTransaction) bool {
	// Validate transaction with Synnergy Consensus
	isValid := automation.consensusEngine.ValidateCrossShardTransaction(transaction)
	return isValid
}

// executeAndLogTransaction executes the validated transaction and logs it in the ledger
func (automation *CrossShardTransactionExecutionAutomation) executeAndLogTransaction(transaction sharding.CrossShardTransaction) {
	// Execute the transaction between shards
	err := automation.shardManager.ExecuteCrossShardTransaction(transaction)
	if err != nil {
		fmt.Printf("Error executing cross-shard transaction %s: %v\n", transaction.ID, err)
		return
	}

	// Log the transaction execution into the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("cross-shard-transaction-%s", transaction.ID),
		Timestamp: time.Now().Unix(),
		Type:      "Cross-Shard Transaction",
		Status:    "Completed",
		Details:   fmt.Sprintf("Cross-shard transaction %s successfully executed.", transaction.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log cross-shard transaction %s: %v\n", transaction.ID, err)
	} else {
		fmt.Printf("Cross-shard transaction %s executed and logged successfully.\n", transaction.ID)
	}
}

// encryptData encrypts sensitive transaction details before storing them in the ledger
func (automation *CrossShardTransactionExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
