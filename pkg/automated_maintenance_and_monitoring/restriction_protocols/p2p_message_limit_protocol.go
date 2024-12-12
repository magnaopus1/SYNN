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
	P2PMessageCheckInterval      = 5 * time.Second // Interval for checking P2P messages
	MaxAllowedP2PMessagesPerNode = 100             // Maximum allowed P2P messages per node per interval
)

// P2PMessageLimitRestrictionAutomation monitors and restricts P2P message limits across the network
type P2PMessageLimitRestrictionAutomation struct {
	consensusSystem     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	stateMutex          *sync.RWMutex
	messageCountPerNode map[string]int // Tracks the number of P2P messages per node
}

// NewP2PMessageLimitRestrictionAutomation initializes and returns an instance of P2PMessageLimitRestrictionAutomation
func NewP2PMessageLimitRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *P2PMessageLimitRestrictionAutomation {
	return &P2PMessageLimitRestrictionAutomation{
		consensusSystem:     consensusSystem,
		ledgerInstance:      ledgerInstance,
		stateMutex:          stateMutex,
		messageCountPerNode: make(map[string]int),
	}
}

// StartP2PMessageMonitoring starts continuous monitoring of P2P messages across the network
func (automation *P2PMessageLimitRestrictionAutomation) StartP2PMessageMonitoring() {
	ticker := time.NewTicker(P2PMessageCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorP2PMessages()
		}
	}()
}

// monitorP2PMessages checks for excessive P2P messages from nodes and enforces restrictions if necessary
func (automation *P2PMessageLimitRestrictionAutomation) monitorP2PMessages() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch P2P message count data from Synnergy Consensus
	messageData := automation.consensusSystem.GetP2PMessageData()

	for nodeID, messageCount := range messageData {
		// Check if the node has exceeded the allowed number of P2P messages
		if messageCount > MaxAllowedP2PMessagesPerNode {
			automation.flagP2PMessageViolation(nodeID, messageCount, "Exceeded allowed P2P message limit")
		}
	}
}

// flagP2PMessageViolation flags a node's P2P message violation and logs it in the ledger
func (automation *P2PMessageLimitRestrictionAutomation) flagP2PMessageViolation(nodeID string, messageCount int, reason string) {
	fmt.Printf("P2P message violation: Node ID %s, Message Count: %d, Reason: %s\n", nodeID, messageCount, reason)

	// Log the violation in the ledger
	automation.logP2PMessageViolation(nodeID, messageCount, reason)
}

// logP2PMessageViolation logs the flagged P2P message violation into the ledger with full details
func (automation *P2PMessageLimitRestrictionAutomation) logP2PMessageViolation(nodeID string, messageCount int, violationReason string) {
	// Create a ledger entry for P2P message violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("p2p-message-violation-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "P2P Message Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Node %s exceeded P2P message limit. Message Count: %d. Reason: %s", nodeID, messageCount, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptP2PMessageData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log P2P message violation:", err)
	} else {
		fmt.Println("P2P message violation logged.")
	}
}

// encryptP2PMessageData encrypts the P2P message data before logging for security
func (automation *P2PMessageLimitRestrictionAutomation) encryptP2PMessageData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting P2P message data:", err)
		return data
	}
	return string(encryptedData)
}
