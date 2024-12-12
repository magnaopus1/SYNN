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
	PrivacyCheckInterval         = 30 * time.Second // Interval for checking privacy compliance
	MaxAllowedPrivacyViolations  = 5                // Maximum allowed privacy violations per user or node
)

// PrivacyEnforcementRestrictionAutomation monitors and enforces privacy compliance across the network
type PrivacyEnforcementRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	privacyViolationCount  map[string]int // Tracks privacy violation counts per user or node
}

// NewPrivacyEnforcementRestrictionAutomation initializes and returns an instance of PrivacyEnforcementRestrictionAutomation
func NewPrivacyEnforcementRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *PrivacyEnforcementRestrictionAutomation {
	return &PrivacyEnforcementRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		privacyViolationCount: make(map[string]int),
	}
}

// StartPrivacyMonitoring starts continuous monitoring of privacy compliance
func (automation *PrivacyEnforcementRestrictionAutomation) StartPrivacyMonitoring() {
	ticker := time.NewTicker(PrivacyCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorPrivacyCompliance()
		}
	}()
}

// monitorPrivacyCompliance checks for privacy violations and enforces restrictions if necessary
func (automation *PrivacyEnforcementRestrictionAutomation) monitorPrivacyCompliance() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch privacy violation data from Synnergy Consensus
	privacyData := automation.consensusSystem.GetPrivacyViolationData()

	for entityID, violationCount := range privacyData {
		// Check if the user or node has exceeded the allowed number of privacy violations
		if automation.privacyViolationCount[entityID] > MaxAllowedPrivacyViolations {
			automation.flagPrivacyViolation(entityID, violationCount, "Exceeded allowed privacy violations")
		}
	}
}

// flagPrivacyViolation flags a privacy violation and logs it in the ledger
func (automation *PrivacyEnforcementRestrictionAutomation) flagPrivacyViolation(entityID string, violationCount int, reason string) {
	fmt.Printf("Privacy violation: Entity ID %s, Reason: %s\n", entityID, reason)

	// Log the violation in the ledger
	automation.logPrivacyViolation(entityID, violationCount, reason)
}

// logPrivacyViolation logs the flagged privacy violation into the ledger with full details
func (automation *PrivacyEnforcementRestrictionAutomation) logPrivacyViolation(entityID string, violationCount int, violationReason string) {
	// Create a ledger entry for privacy violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("privacy-violation-%s-%d", entityID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Privacy Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Entity %s violated privacy rules. Violation Count: %d. Reason: %s", entityID, violationCount, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptPrivacyData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log privacy violation:", err)
	} else {
		fmt.Println("Privacy violation logged.")
	}
}

// encryptPrivacyData encrypts the privacy violation data before logging for security
func (automation *PrivacyEnforcementRestrictionAutomation) encryptPrivacyData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting privacy violation data:", err)
		return data
	}
	return string(encryptedData)
}
