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
	FundMovementCheckInterval  = 30 * time.Second // Interval for checking restricted auto fund movements
	MaxAllowedMovementViolations = 5              // Maximum allowed fund movement violations before restriction
)

// RestrictedAutoFundMovementAutomation enforces restrictions on automated fund transfers
type RestrictedAutoFundMovementAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	movementViolationCount   map[string]int // Tracks automated fund movement violations per user or entity
}

// NewRestrictedAutoFundMovementAutomation initializes and returns an instance of RestrictedAutoFundMovementAutomation
func NewRestrictedAutoFundMovementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedAutoFundMovementAutomation {
	return &RestrictedAutoFundMovementAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		movementViolationCount: make(map[string]int),
	}
}

// StartFundMovementMonitoring starts continuous monitoring of restricted automated fund movement activity
func (automation *RestrictedAutoFundMovementAutomation) StartFundMovementMonitoring() {
	ticker := time.NewTicker(FundMovementCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorFundMovements()
		}
	}()
}

// monitorFundMovements checks for restricted fund movement violations and enforces restrictions if necessary
func (automation *RestrictedAutoFundMovementAutomation) monitorFundMovements() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch fund movement data from Synnergy Consensus
	movementData := automation.consensusSystem.GetFundMovementData()

	for userID, movementStatus := range movementData {
		// Check if the user has violated the fund movement rules
		if movementStatus == "violation" {
			automation.flagMovementViolation(userID, "Unauthorized automated fund movement detected")
		}
	}
}

// flagMovementViolation flags a user's automated fund movement violation and logs it in the ledger
func (automation *RestrictedAutoFundMovementAutomation) flagMovementViolation(userID string, reason string) {
	fmt.Printf("Auto fund movement violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.movementViolationCount[userID]++

	// Log the violation in the ledger
	automation.logMovementViolation(userID, reason)

	// Check if the user has exceeded the allowed number of movement violations
	if automation.movementViolationCount[userID] >= MaxAllowedMovementViolations {
		automation.restrictFundMovement(userID)
	}
}

// logMovementViolation logs the flagged fund movement violation into the ledger with full details
func (automation *RestrictedAutoFundMovementAutomation) logMovementViolation(userID string, violationReason string) {
	// Create a ledger entry for fund movement violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("fund-movement-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Fund Movement Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated fund movement restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptMovementData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log fund movement violation:", err)
	} else {
		fmt.Println("Fund movement violation logged.")
	}
}

// restrictFundMovement restricts auto fund movement access for a user after exceeding allowed violations
func (automation *RestrictedAutoFundMovementAutomation) restrictFundMovement(userID string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("fund-movement-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Fund Movement Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from automated fund movements due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptMovementData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log fund movement restriction:", err)
	} else {
		fmt.Println("Fund movement restriction applied.")
	}
}

// encryptMovementData encrypts the fund movement data before logging for security
func (automation *RestrictedAutoFundMovementAutomation) encryptMovementData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting movement data:", err)
		return data
	}
	return string(encryptedData)
}
