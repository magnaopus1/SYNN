package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
)

const (
	WeightedDistributionInterval       = 2 * time.Minute  // Time interval to check and apply weighted distribution
	WeightedDistributionLogKey         = "weightedDistKey" // Encryption key for logging distribution events
	MaxWeightThreshold                 = 100               // Maximum weight threshold for a node
	MinWeightThreshold                 = 1                 // Minimum weight threshold for a node
)

// WeightedDistributionAutomation handles distributing loads across nodes based on weights.
type WeightedDistributionAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging distribution events
	consensusSystem *consensus.SynnergyConsensus // Reference to the consensus system for transaction handling
	syncMutex       *sync.RWMutex                // Mutex to handle concurrent distribution
}

// NewWeightedDistributionAutomation creates a new instance of the automation process.
func NewWeightedDistributionAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, syncMutex *sync.RWMutex) *WeightedDistributionAutomation {
	return &WeightedDistributionAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		syncMutex:       syncMutex,
	}
}

// StartWeightedDistributionAutomation starts the continuous loop for distributing loads based on node weight.
func (automation *WeightedDistributionAutomation) StartWeightedDistributionAutomation() {
	ticker := time.NewTicker(WeightedDistributionInterval)
	go func() {
		for range ticker.C {
			automation.distributeLoadByWeight()
		}
	}()
}

// distributeLoadByWeight fetches active nodes and distributes loads based on their weights.
func (automation *WeightedDistributionAutomation) distributeLoadByWeight() {
	automation.syncMutex.Lock()
	defer automation.syncMutex.Unlock()

	// Fetch active nodes from the consensus system
	nodes, err := automation.consensusSystem.GetActiveNodes()
	if err != nil {
		fmt.Printf("Error fetching active nodes for weighted distribution: %v\n", err)
		return
	}

	// Calculate total weight to use for distribution
	totalWeight := automation.calculateTotalWeight(nodes)

	// Distribute transactions to nodes based on their weight
	for _, node := range nodes {
		if node.Weight >= MinWeightThreshold && node.Weight <= MaxWeightThreshold {
			automation.distributeToNode(node, totalWeight)
		} else {
			fmt.Printf("Node %s has invalid weight %d, skipping...\n", node.ID, node.Weight)
		}
	}
}

// calculateTotalWeight calculates the sum of weights for all active nodes.
func (automation *WeightedDistributionAutomation) calculateTotalWeight(nodes []common.Node) int {
	totalWeight := 0
	for _, node := range nodes {
		if node.Weight >= MinWeightThreshold && node.Weight <= MaxWeightThreshold {
			totalWeight += node.Weight
		}
	}
	return totalWeight
}

// distributeToNode distributes the load to a specific node based on its weight.
func (automation *WeightedDistributionAutomation) distributeToNode(node common.Node, totalWeight int) {
	if totalWeight == 0 {
		fmt.Printf("Total weight is 0, cannot distribute load to node %s\n", node.ID)
		return
	}

	// Calculate the portion of load for this node
	loadPercentage := float64(node.Weight) / float64(totalWeight)

	// Fetch transactions that need distribution
	transactions, err := automation.consensusSystem.GetPendingTransactions()
	if err != nil {
		fmt.Printf("Error fetching pending transactions: %v\n", err)
		return
	}

	// Calculate the number of transactions to assign to this node
	numTransactions := int(loadPercentage * float64(len(transactions)))

	// Distribute transactions to this node
	for i := 0; i < numTransactions; i++ {
		tx := transactions[i]
		if err := automation.consensusSystem.AssignTransactionToNode(tx, node); err != nil {
			fmt.Printf("Error assigning transaction %s to node %s: %v\n", tx.ID, node.ID, err)
			continue
		}
		automation.logDistributionEvent(tx, node)
	}
}

// logDistributionEvent logs the distribution event to the ledger with encrypted details.
func (automation *WeightedDistributionAutomation) logDistributionEvent(tx common.Transaction, node common.Node) {
	// Encrypt the log entry for security
	encryptedLogDetails, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Transaction %s assigned to node %s with weight %d", tx.ID, node.ID, node.Weight)), WeightedDistributionLogKey)
	if err != nil {
		fmt.Printf("Error encrypting distribution log for transaction %s: %v\n", tx.ID, err)
		return
	}

	// Create a ledger entry for the weighted distribution event
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("TX-DISTRIBUTION-%s-%d", tx.ID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Distribution",
		Status:    "Success",
		Details:   string(encryptedLogDetails),
	}

	// Log the distribution event in the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log distribution event for transaction %s: %v\n", tx.ID, err)
	}
}
