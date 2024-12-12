package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/blocks"
	"synnergy_network_demo/network"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
)

const (
	ThroughputMonitorInterval  = 5 * time.Second // Interval to monitor transaction throughput
	ThroughputThreshold        = 10000           // Threshold for triggering an alert
)

// TransactionThroughputAutomation handles monitoring and optimizing transaction throughput
type TransactionThroughputAutomation struct {
	networkInstance   *network.Network
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	blockchain        *blocks.Blockchain
	throughputMutex   *sync.RWMutex
}

// NewTransactionThroughputAutomation initializes the transaction throughput automation
func NewTransactionThroughputAutomation(networkInstance *network.Network, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, blockchain *blocks.Blockchain, throughputMutex *sync.RWMutex) *TransactionThroughputAutomation {
	return &TransactionThroughputAutomation{
		networkInstance: networkInstance,
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		blockchain:      blockchain,
		throughputMutex: throughputMutex,
	}
}

// StartTransactionThroughputMonitor starts monitoring transaction throughput in a continuous loop
func (automation *TransactionThroughputAutomation) StartTransactionThroughputMonitor() {
	ticker := time.NewTicker(ThroughputMonitorInterval)

	go func() {
		for range ticker.C {
			automation.monitorThroughput()
		}
	}()
}

// monitorThroughput checks transaction throughput and applies adjustments as needed
func (automation *TransactionThroughputAutomation) monitorThroughput() {
	automation.throughputMutex.Lock()
	defer automation.throughputMutex.Unlock()

	// Get current transaction throughput
	currentThroughput := automation.consensusEngine.GetTransactionThroughput()

	// If throughput exceeds threshold, trigger alert and optimize
	if currentThroughput > ThroughputThreshold {
		fmt.Printf("Transaction throughput exceeded: %d transactions per second.\n", currentThroughput)
		automation.optimizeThroughput()
		automation.logThroughputEvent(currentThroughput)
	}
}

// optimizeThroughput optimizes the network and consensus system to handle high transaction loads
func (automation *TransactionThroughputAutomation) optimizeThroughput() {
	fmt.Println("Optimizing transaction throughput to handle high loads...")

	// Apply consensus engine optimizations
	err := automation.consensusEngine.OptimizeForHighLoad()
	if err != nil {
		fmt.Printf("Error optimizing consensus engine for high load: %v\n", err)
	}

	// Notify network to adjust resources for high throughput
	err = automation.networkInstance.AdjustResourcesForThroughput(ThroughputThreshold)
	if err != nil {
		fmt.Printf("Error adjusting network resources: %v\n", err)
	}
}

// logThroughputEvent logs the transaction throughput event in the ledger for auditing
func (automation *TransactionThroughputAutomation) logThroughputEvent(throughput int) {
	entryDetails := fmt.Sprintf("Transaction throughput exceeded threshold: %d transactions per second.", throughput)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := ledger.LedgerEntry{
		ID:        fmt.Sprintf("transaction-throughput-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Throughput Alert",
		Status:    "Handled",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log throughput event in the ledger: %v\n", err)
	} else {
		fmt.Println("Throughput event successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *TransactionThroughputAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualThroughputOptimization allows administrators to manually optimize the system for high throughput
func (automation *TransactionThroughputAutomation) TriggerManualThroughputOptimization() {
	fmt.Println("Manually optimizing the system for high transaction throughput...")
	automation.optimizeThroughput()
}
