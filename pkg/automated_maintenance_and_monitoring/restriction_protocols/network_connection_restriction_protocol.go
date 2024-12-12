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
	NetworkConnectionCheckInterval = 10 * time.Second // Interval for checking network connection issues
	MaxAllowedFailedConnections    = 3                // Maximum allowed failed connection attempts before restriction
)

// NetworkConnectionRestrictionAutomation monitors and restricts network connection issues
type NetworkConnectionRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	failedConnectionCount  map[string]int // Tracks failed connection attempts per node
}

// NewNetworkConnectionRestrictionAutomation initializes and returns an instance of NetworkConnectionRestrictionAutomation
func NewNetworkConnectionRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NetworkConnectionRestrictionAutomation {
	return &NetworkConnectionRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		failedConnectionCount: make(map[string]int),
	}
}

// StartConnectionMonitoring starts continuous monitoring of network connection issues
func (automation *NetworkConnectionRestrictionAutomation) StartConnectionMonitoring() {
	ticker := time.NewTicker(NetworkConnectionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorNetworkConnections()
		}
	}()
}

// monitorNetworkConnections checks for failed network connections and enforces restrictions if necessary
func (automation *NetworkConnectionRestrictionAutomation) monitorNetworkConnections() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch failed connection data from Synnergy Consensus
	connectionData := automation.consensusSystem.GetFailedConnectionData()

	for nodeID, failedAttempts := range connectionData {
		// Check if the node has exceeded the allowed number of failed connections
		if automation.failedConnectionCount[nodeID] > MaxAllowedFailedConnections {
			automation.flagConnectionViolation(nodeID, failedAttempts, "Exceeded allowed failed connection attempts")
		}
	}
}

// flagConnectionViolation flags a node's connection issue and logs it in the ledger
func (automation *NetworkConnectionRestrictionAutomation) flagConnectionViolation(nodeID string, failedAttempts int, reason string) {
	fmt.Printf("Network connection violation: Node ID %s, Failed Attempts: %d, Reason: %s\n", nodeID, failedAttempts, reason)

	// Log the violation in the ledger
	automation.logConnectionViolation(nodeID, failedAttempts, reason)
}

// logConnectionViolation logs the flagged network connection violation into the ledger with full details
func (automation *NetworkConnectionRestrictionAutomation) logConnectionViolation(nodeID string, failedAttempts int, violationReason string) {
	// Create a ledger entry for network connection violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("network-connection-violation-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Network Connection Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Node %s violated network connection rules. Failed Attempts: %d. Reason: %s", nodeID, failedAttempts, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptConnectionData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log network connection violation:", err)
	} else {
		fmt.Println("Network connection violation logged.")
	}
}

// encryptConnectionData encrypts the connection data before logging for security
func (automation *NetworkConnectionRestrictionAutomation) encryptConnectionData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting connection data:", err)
		return data
	}
	return string(encryptedData)
}
