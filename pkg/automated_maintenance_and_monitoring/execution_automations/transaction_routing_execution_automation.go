package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/transactions"
	"synnergy_network_demo/blocks"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/network"
)

const (
	TransactionRoutingInterval = 10 * time.Second // Time interval for routing transactions
	MaxPendingTransactions     = 500              // Max transactions before routing is triggered
)

// TransactionRoutingAutomation handles the routing of pending transactions across the network
type TransactionRoutingAutomation struct {
	networkInstance   *network.Network
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	blockchain        *blocks.Blockchain
	routingMutex      *sync.RWMutex
}

// NewTransactionRoutingAutomation initializes the transaction routing automation
func NewTransactionRoutingAutomation(networkInstance *network.Network, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, blockchain *blocks.Blockchain, routingMutex *sync.RWMutex) *TransactionRoutingAutomation {
	return &TransactionRoutingAutomation{
		networkInstance: networkInstance,
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		blockchain:      blockchain,
		routingMutex:    routingMutex,
	}
}

// StartTransactionRoutingMonitor starts a continuous loop that monitors and routes pending transactions
func (automation *TransactionRoutingAutomation) StartTransactionRoutingMonitor() {
	ticker := time.NewTicker(TransactionRoutingInterval)

	go func() {
		for range ticker.C {
			automation.routePendingTransactions()
		}
	}()
}

// routePendingTransactions checks for pending transactions and routes them across the network
func (automation *TransactionRoutingAutomation) routePendingTransactions() {
	automation.routingMutex.Lock()
	defer automation.routingMutex.Unlock()

	// Get the pending transactions
	pendingTransactions := automation.blockchain.GetPendingTransactions()

	// If the number of pending transactions exceeds the limit, route them
	if len(pendingTransactions) >= MaxPendingTransactions {
		err := automation.processAndRouteTransactions(pendingTransactions)
		if err != nil {
			fmt.Printf("Error routing transactions: %v\n", err)
			return
		}

		// Log the transaction routing into the ledger
		automation.logTransactionRouting(pendingTransactions)
	}
}

// processAndRouteTransactions routes the transactions through the Synnergy network consensus
func (automation *TransactionRoutingAutomation) processAndRouteTransactions(transactions []*transactions.Transaction) error {
	// Validate the transactions using Synnergy Consensus
	valid, err := automation.consensusEngine.ValidateTransactions(transactions)
	if err != nil || !valid {
		return fmt.Errorf("consensus validation failed for transactions: %v", err)
	}

	// Broadcast the validated transactions across the network
	err = automation.networkInstance.BroadcastTransactions(transactions)
	if err != nil {
		return fmt.Errorf("failed to broadcast transactions: %v", err)
	}

	fmt.Println("Transactions successfully routed across the network.")
	return nil
}

// logTransactionRouting logs the transaction routing event into the ledger
func (automation *TransactionRoutingAutomation) logTransactionRouting(transactions []*transactions.Transaction) {
	entryDetails := fmt.Sprintf("Routed %d transactions across the network.", len(transactions))
	encryptedDetails := automation.encryptData(entryDetails)

	entry := ledger.LedgerEntry{
		ID:        fmt.Sprintf("transaction-routing-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Routing",
		Status:    "Success",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log transaction routing event in the ledger: %v\n", err)
	} else {
		fmt.Println("Transaction routing successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *TransactionRoutingAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualTransactionRouting allows administrators to manually trigger transaction routing
func (automation *TransactionRoutingAutomation) TriggerManualTransactionRouting() {
	fmt.Println("Manually triggering transaction routing...")

	automation.routePendingTransactions()
}
