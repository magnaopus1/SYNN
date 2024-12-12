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
	OracleDataValidationCheckInterval = 15 * time.Second // Interval for checking oracle data validation
	MaxAllowedInvalidOracleData       = 5                // Maximum allowed invalid oracle data entries
)

// OracleDataValidationRestrictionAutomation monitors and restricts invalid oracle data entries
type OracleDataValidationRestrictionAutomation struct {
	consensusSystem         *consensus.SynnergyConsensus
	ledgerInstance          *ledger.Ledger
	stateMutex              *sync.RWMutex
	invalidOracleDataCount  map[string]int // Tracks invalid oracle data counts per oracle
}

// NewOracleDataValidationRestrictionAutomation initializes and returns an instance of OracleDataValidationRestrictionAutomation
func NewOracleDataValidationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *OracleDataValidationRestrictionAutomation {
	return &OracleDataValidationRestrictionAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		invalidOracleDataCount: make(map[string]int),
	}
}

// StartOracleDataMonitoring starts continuous monitoring of oracle data validation
func (automation *OracleDataValidationRestrictionAutomation) StartOracleDataMonitoring() {
	ticker := time.NewTicker(OracleDataValidationCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorOracleData()
		}
	}()
}

// monitorOracleData checks for invalid oracle data entries and enforces restrictions if necessary
func (automation *OracleDataValidationRestrictionAutomation) monitorOracleData() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch oracle data validation results from Synnergy Consensus
	oracleData := automation.consensusSystem.GetOracleDataValidationResults()

	for oracleID, validationResult := range oracleData {
		// Check if the oracle has provided invalid data more than the allowed limit
		if automation.invalidOracleDataCount[oracleID] > MaxAllowedInvalidOracleData {
			automation.flagOracleDataViolation(oracleID, validationResult, "Exceeded allowed invalid oracle data entries")
		}
	}
}

// flagOracleDataViolation flags an oracle's invalid data entry violation and logs it in the ledger
func (automation *OracleDataValidationRestrictionAutomation) flagOracleDataViolation(oracleID string, validationResult string, reason string) {
	fmt.Printf("Oracle data validation violation: Oracle ID %s, Reason: %s\n", oracleID, reason)

	// Log the violation in the ledger
	automation.logOracleDataViolation(oracleID, validationResult, reason)
}

// logOracleDataViolation logs the flagged oracle data validation violation into the ledger with full details
func (automation *OracleDataValidationRestrictionAutomation) logOracleDataViolation(oracleID string, validationResult string, violationReason string) {
	// Create a ledger entry for oracle data validation violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("oracle-data-violation-%s-%d", oracleID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Oracle Data Validation Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Oracle %s provided invalid data. Validation Result: %s. Reason: %s", oracleID, validationResult, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptOracleData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log oracle data validation violation:", err)
	} else {
		fmt.Println("Oracle data validation violation logged.")
	}
}

// encryptOracleData encrypts the oracle data before logging for security
func (automation *OracleDataValidationRestrictionAutomation) encryptOracleData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting oracle data:", err)
		return data
	}
	return string(encryptedData)
}
