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
	SecurityCheckInterval     = 1 * time.Second  // Interval for checking security breaches
	MaxSecurityBreachAttempts = 3                // Maximum allowed security breach attempts
)

// SecurityBreachResponseAutomation monitors and responds to network security breaches
type SecurityBreachResponseAutomation struct {
	consensusSystem         *consensus.SynnergyConsensus
	ledgerInstance          *ledger.Ledger
	stateMutex              *sync.RWMutex
	securityBreachAttempts  map[string]int // Tracks breach attempts by nodes or users
}

// NewSecurityBreachResponseAutomation initializes the SecurityBreachResponseAutomation struct
func NewSecurityBreachResponseAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SecurityBreachResponseAutomation {
	return &SecurityBreachResponseAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		securityBreachAttempts: make(map[string]int),
	}
}

// StartSecurityBreachMonitoring starts the continuous security breach detection process
func (automation *SecurityBreachResponseAutomation) StartSecurityBreachMonitoring() {
	ticker := time.NewTicker(SecurityCheckInterval)

	go func() {
		for range ticker.C {
			automation.detectAndRespondToSecurityBreaches()
		}
	}()
}

// detectAndRespondToSecurityBreaches monitors and responds to any detected security breaches
func (automation *SecurityBreachResponseAutomation) detectAndRespondToSecurityBreaches() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch breach attempt data from Synnergy Consensus
	breachData := automation.consensusSystem.GetSecurityBreachAttempts()

	for nodeID, breachStatus := range breachData {
		if breachStatus == "breach" {
			automation.logSecurityBreach(nodeID, "Security breach detected.")
			automation.increaseBreachAttemptCount(nodeID)

			// Check if the node has exceeded the allowed number of breach attempts
			if automation.securityBreachAttempts[nodeID] >= MaxSecurityBreachAttempts {
				automation.restrictNodeAccess(nodeID)
			}
		}
	}
}

// logSecurityBreach logs the security breach to the ledger
func (automation *SecurityBreachResponseAutomation) logSecurityBreach(nodeID string, breachDetails string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("security-breach-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Security Breach",
		Status:    "Detected",
		Details:   fmt.Sprintf("Node %s detected security breach: %s", nodeID, breachDetails),
	}

	// Encrypt the security breach data before logging
	encryptedDetails := automation.encryptSecurityBreachData(entry.Details)
	entry.Details = encryptedDetails

	// Add the breach entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log security breach:", err)
	} else {
		fmt.Println("Security breach logged successfully.")
	}
}

// increaseBreachAttemptCount increases the breach attempt count for the specific node
func (automation *SecurityBreachResponseAutomation) increaseBreachAttemptCount(nodeID string) {
	automation.securityBreachAttempts[nodeID]++
	fmt.Printf("Breach attempt count for node %s increased to %d.\n", nodeID, automation.securityBreachAttempts[nodeID])
}

// restrictNodeAccess restricts a node's access after multiple security breach attempts
func (automation *SecurityBreachResponseAutomation) restrictNodeAccess(nodeID string) {
	fmt.Printf("Node %s has exceeded the allowed number of breach attempts and is being restricted.\n", nodeID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("security-restriction-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Security Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Node %s restricted from the network after repeated security breaches.", nodeID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptSecurityBreachData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log security restriction:", err)
	} else {
		fmt.Println("Security restriction applied successfully.")
	}

	// Remove node from consensus participation
	automation.consensusSystem.RestrictNode(nodeID)
}

// encryptSecurityBreachData encrypts sensitive data before logging it in the ledger
func (automation *SecurityBreachResponseAutomation) encryptSecurityBreachData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting security breach data:", err)
		return data
	}
	return string(encryptedData)
}
