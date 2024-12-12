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
	SensitiveDataAccessCheckInterval = 1 * time.Second // Interval for checking sensitive data access
	MaxSensitiveAccessAttempts       = 3               // Maximum allowed attempts for accessing sensitive data
)

// SensitiveDataAccessRestrictionAutomation handles restricting access to sensitive data based on violation thresholds
type SensitiveDataAccessRestrictionAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	sensitiveAccessAttempts  map[string]int // Tracks access attempts for sensitive data by users or nodes
}

// NewSensitiveDataAccessRestrictionAutomation initializes the SensitiveDataAccessRestrictionAutomation struct
func NewSensitiveDataAccessRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SensitiveDataAccessRestrictionAutomation {
	return &SensitiveDataAccessRestrictionAutomation{
		consensusSystem:         consensusSystem,
		ledgerInstance:          ledgerInstance,
		stateMutex:              stateMutex,
		sensitiveAccessAttempts: make(map[string]int),
	}
}

// StartSensitiveDataAccessMonitoring starts the continuous monitoring process for sensitive data access
func (automation *SensitiveDataAccessRestrictionAutomation) StartSensitiveDataAccessMonitoring() {
	ticker := time.NewTicker(SensitiveDataAccessCheckInterval)

	go func() {
		for range ticker.C {
			automation.detectAndRestrictSensitiveDataAccess()
		}
	}()
}

// detectAndRestrictSensitiveDataAccess monitors sensitive data access attempts and enforces restrictions if necessary
func (automation *SensitiveDataAccessRestrictionAutomation) detectAndRestrictSensitiveDataAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch sensitive data access attempts from the consensus system
	accessData := automation.consensusSystem.GetSensitiveDataAccessAttempts()

	for nodeID, accessStatus := range accessData {
		if accessStatus == "attempted" {
			automation.logSensitiveAccess(nodeID, "Unauthorized attempt to access sensitive data detected.")
			automation.increaseSensitiveAccessAttemptCount(nodeID)

			// If a node exceeds the allowed number of sensitive data access attempts, restrict access
			if automation.sensitiveAccessAttempts[nodeID] >= MaxSensitiveAccessAttempts {
				automation.restrictSensitiveDataAccess(nodeID)
			}
		}
	}
}

// logSensitiveAccess logs sensitive data access attempts in the ledger
func (automation *SensitiveDataAccessRestrictionAutomation) logSensitiveAccess(nodeID string, accessDetails string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("sensitive-access-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Sensitive Data Access Attempt",
		Status:    "Attempted",
		Details:   fmt.Sprintf("Node %s attempted to access sensitive data: %s", nodeID, accessDetails),
	}

	// Encrypt the sensitive access data before logging
	encryptedDetails := automation.encryptSensitiveData(entry.Details)
	entry.Details = encryptedDetails

	// Add the sensitive access entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log sensitive data access attempt:", err)
	} else {
		fmt.Println("Sensitive data access attempt logged successfully.")
	}
}

// increaseSensitiveAccessAttemptCount increases the count of access attempts for a specific node
func (automation *SensitiveDataAccessRestrictionAutomation) increaseSensitiveAccessAttemptCount(nodeID string) {
	automation.sensitiveAccessAttempts[nodeID]++
	fmt.Printf("Sensitive data access attempt count for node %s increased to %d.\n", nodeID, automation.sensitiveAccessAttempts[nodeID])
}

// restrictSensitiveDataAccess restricts access for a node that has made too many sensitive data access attempts
func (automation *SensitiveDataAccessRestrictionAutomation) restrictSensitiveDataAccess(nodeID string) {
	fmt.Printf("Node %s has exceeded sensitive data access attempts and is being restricted.\n", nodeID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("sensitive-access-restriction-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Sensitive Data Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Node %s restricted from accessing sensitive data after repeated attempts.", nodeID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptSensitiveData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log sensitive data access restriction:", err)
	} else {
		fmt.Println("Sensitive data access restriction applied successfully.")
	}

	// Remove node from sensitive data access within consensus
	automation.consensusSystem.RestrictNodeAccessToSensitiveData(nodeID)
}

// encryptSensitiveData encrypts sensitive data access details before storing them in the ledger
func (automation *SensitiveDataAccessRestrictionAutomation) encryptSensitiveData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting sensitive data access details:", err)
		return data
	}
	return string(encryptedData)
}
