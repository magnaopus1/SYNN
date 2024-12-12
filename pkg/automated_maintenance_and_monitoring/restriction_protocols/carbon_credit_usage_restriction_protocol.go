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
	CarbonCreditCheckInterval    = 15 * time.Second // Interval for monitoring carbon credit usage
	MaxCarbonCreditsPerPeriod    = 1000.0           // Maximum carbon credits usage allowed per time period
	CarbonCreditTimePeriod       = 30 * 24 * time.Hour // Time period in which max usage is calculated (30 days)
)

// CarbonCreditUsageRestrictionAutomation monitors and restricts carbon credit usage across the network
type CarbonCreditUsageRestrictionAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	carbonCreditUsage     map[string]float64 // Tracks carbon credit usage by account
}

// NewCarbonCreditUsageRestrictionAutomation initializes and returns an instance of CarbonCreditUsageRestrictionAutomation
func NewCarbonCreditUsageRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CarbonCreditUsageRestrictionAutomation {
	return &CarbonCreditUsageRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		stateMutex:        stateMutex,
		carbonCreditUsage: make(map[string]float64),
	}
}

// StartCarbonCreditMonitoring starts continuous monitoring of carbon credit usage for restrictions
func (automation *CarbonCreditUsageRestrictionAutomation) StartCarbonCreditMonitoring() {
	ticker := time.NewTicker(CarbonCreditCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorCarbonCreditUsage()
		}
	}()
}

// monitorCarbonCreditUsage checks carbon credit usage across accounts and enforces usage limits
func (automation *CarbonCreditUsageRestrictionAutomation) monitorCarbonCreditUsage() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent carbon credit usage records from Synnergy Consensus
	recentCreditUsages := automation.consensusSystem.GetRecentCarbonCreditUsages()

	for _, usage := range recentCreditUsages {
		// Validate if the usage exceeds the maximum allowed for the time period
		if !automation.validateCarbonCreditUsage(usage) {
			automation.flagCarbonCreditViolation(usage, "Carbon credit usage exceeded the maximum limit")
		}
	}
}

// validateCarbonCreditUsage checks if the carbon credit usage exceeds the allowable limit for the period
func (automation *CarbonCreditUsageRestrictionAutomation) validateCarbonCreditUsage(usage common.CarbonCreditUsage) bool {
	currentUsage := automation.carbonCreditUsage[usage.AccountID]
	if currentUsage+usage.CreditsUsed > MaxCarbonCreditsPerPeriod {
		return false
	}

	// Update the usage for the account
	automation.carbonCreditUsage[usage.AccountID] += usage.CreditsUsed
	return true
}

// flagCarbonCreditViolation flags an account that violates carbon credit usage rules and logs it in the ledger
func (automation *CarbonCreditUsageRestrictionAutomation) flagCarbonCreditViolation(usage common.CarbonCreditUsage, reason string) {
	fmt.Printf("Carbon credit usage violation: Account %s, Reason: %s\n", usage.AccountID, reason)

	// Log the violation into the ledger
	automation.logCarbonCreditViolation(usage, reason)
}

// logCarbonCreditViolation logs the flagged carbon credit violation into the ledger with full details
func (automation *CarbonCreditUsageRestrictionAutomation) logCarbonCreditViolation(usage common.CarbonCreditUsage, violationReason string) {
	// Encrypt the carbon credit violation data
	encryptedData := automation.encryptCarbonCreditUsageData(usage)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("carbon-credit-violation-%s-%d", usage.AccountID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Carbon Credit Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Account %s flagged for carbon credit usage violation. Reason: %s. Encrypted Data: %s", usage.AccountID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log carbon credit usage violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Carbon credit usage violation logged for account: %s\n", usage.AccountID)
	}
}

// encryptCarbonCreditUsageData encrypts carbon credit usage data before logging
func (automation *CarbonCreditUsageRestrictionAutomation) encryptCarbonCreditUsageData(usage common.CarbonCreditUsage) string {
	data := fmt.Sprintf("Account: %s, Credits Used: %.2f", usage.AccountID, usage.CreditsUsed)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting carbon credit usage data:", err)
		return data
	}
	return string(encryptedData)
}
