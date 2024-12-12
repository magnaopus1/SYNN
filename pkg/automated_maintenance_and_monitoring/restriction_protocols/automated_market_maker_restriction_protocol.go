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
	AMMRestrictionCheckInterval = 10 * time.Second  // Interval for monitoring AMM operations
	MaxLiquidityWithdrawals     = 1000.0            // Maximum allowable withdrawal of liquidity before restriction
	MinTradeAmount              = 0.01              // Minimum trade amount to prevent market manipulation
)

// AMMRestrictionAutomation handles the restrictions and monitoring of Automated Market Maker operations
type AMMRestrictionAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	flaggedAMMOperations  map[string]int // Tracks flagged AMM operations
}

// NewAMMRestrictionAutomation initializes and returns an instance of AMMRestrictionAutomation
func NewAMMRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AMMRestrictionAutomation {
	return &AMMRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		flaggedAMMOperations: make(map[string]int),
	}
}

// StartAMMMonitoring starts continuous monitoring of Automated Market Maker operations for restrictions
func (automation *AMMRestrictionAutomation) StartAMMMonitoring() {
	ticker := time.NewTicker(AMMRestrictionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorAMMOperations()
		}
	}()
}

// monitorAMMOperations checks AMM activity for compliance with restrictions and flags violations
func (automation *AMMRestrictionAutomation) monitorAMMOperations() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent AMM operations from Synnergy Consensus
	recentOperations := automation.consensusSystem.GetRecentAMMOperations()

	for _, operation := range recentOperations {
		// Validate liquidity withdrawals and trade amounts
		if !automation.validateLiquidityWithdrawal(operation) {
			automation.flagAMMOperation(operation, "Liquidity withdrawal exceeds the maximum allowed")
			continue
		}

		if !automation.validateTradeAmount(operation) {
			automation.flagAMMOperation(operation, "Trade amount below the minimum allowed")
		}
	}
}

// validateLiquidityWithdrawal checks if the liquidity withdrawal exceeds the allowable maximum
func (automation *AMMRestrictionAutomation) validateLiquidityWithdrawal(operation common.AMMOperation) bool {
	return operation.LiquidityWithdrawal <= MaxLiquidityWithdrawals
}

// validateTradeAmount checks if the trade amount meets the minimum allowable amount to prevent manipulation
func (automation *AMMRestrictionAutomation) validateTradeAmount(operation common.AMMOperation) bool {
	return operation.TradeAmount >= MinTradeAmount
}

// flagAMMOperation flags an AMM operation that violates restrictions and logs it in the ledger
func (automation *AMMRestrictionAutomation) flagAMMOperation(operation common.AMMOperation, reason string) {
	fmt.Printf("AMM operation flagged: %s, Reason: %s\n", operation.OperationID, reason)

	// Track flagged AMM operations
	automation.flaggedAMMOperations[operation.OperationID]++

	// Log the violation into the ledger
	automation.logAMMViolation(operation, reason)
}

// logAMMViolation logs the flagged AMM operation into the ledger with full details
func (automation *AMMRestrictionAutomation) logAMMViolation(operation common.AMMOperation, violationReason string) {
	// Encrypt the operation data
	encryptedOperationData := automation.encryptAMMOperationData(operation)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("amm-violation-%s-%d", operation.OperationID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "AMM Operation Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("AMM operation (%s) flagged for violation. Reason: %s. Encrypted Data: %s", operation.OperationID, violationReason, encryptedOperationData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log AMM operation violation into ledger: %v\n", err)
	} else {
		fmt.Printf("AMM operation violation logged for operation: %s\n", operation.OperationID)
	}
}

// encryptAMMOperationData encrypts AMM operation data before logging
func (automation *AMMRestrictionAutomation) encryptAMMOperationData(operation common.AMMOperation) string {
	data := fmt.Sprintf("Operation ID: %s, Liquidity Withdrawal: %.2f, Trade Amount: %.2f", operation.OperationID, operation.LiquidityWithdrawal, operation.TradeAmount)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting AMM operation data:", err)
		return data
	}
	return string(encryptedData)
}
