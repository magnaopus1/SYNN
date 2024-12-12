package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/dex"
)

// Configuration for DEX enforcement automation
const (
	DEXCheckInterval              = 15 * time.Second // Interval to check DEX compliance
	MaxTransactionFrequency       = 100              // Maximum number of transactions allowed per minute
	MinLiquidityThreshold         = 5000             // Minimum liquidity required in a pool
	MaxLiquidityWithdrawals       = 3                // Maximum allowed withdrawals per hour per pool
)

// DEXEnforcementAutomation monitors and enforces compliance for DEX operations
type DEXEnforcementAutomation struct {
	dexManager        *dex.DEXManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	transactionCount  map[string]int // Tracks transaction frequency per user within a time period
	withdrawalCount   map[string]int // Tracks withdrawal frequency per liquidity pool
}

// NewDEXEnforcementAutomation initializes the DEX enforcement automation
func NewDEXEnforcementAutomation(dexManager *dex.DEXManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DEXEnforcementAutomation {
	return &DEXEnforcementAutomation{
		dexManager:       dexManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		transactionCount: make(map[string]int),
		withdrawalCount:  make(map[string]int),
	}
}

// StartDEXEnforcement begins continuous monitoring and enforcement of DEX compliance
func (automation *DEXEnforcementAutomation) StartDEXEnforcement() {
	ticker := time.NewTicker(DEXCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkDEXCompliance()
		}
	}()
}

// checkDEXCompliance monitors DEX transaction frequency, liquidity levels, and withdrawal activity for compliance
func (automation *DEXEnforcementAutomation) checkDEXCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyTransactionFrequency()
	automation.checkLiquidityThresholds()
	automation.enforceWithdrawalLimits()
}

// verifyTransactionFrequency checks if transaction frequency per user exceeds allowable limits
func (automation *DEXEnforcementAutomation) verifyTransactionFrequency() {
	for _, userID := range automation.dexManager.GetActiveUsers() {
		txCount := automation.dexManager.GetUserTransactionCount(userID)

		if txCount > MaxTransactionFrequency {
			fmt.Printf("Transaction frequency violation for user %s.\n", userID)
			automation.handleTransactionViolation(userID)
		}
	}
}

// checkLiquidityThresholds ensures liquidity pools meet minimum liquidity requirements
func (automation *DEXEnforcementAutomation) checkLiquidityThresholds() {
	for _, poolID := range automation.dexManager.GetLiquidityPools() {
		liquidity := automation.dexManager.GetPoolLiquidity(poolID)

		if liquidity < MinLiquidityThreshold {
			fmt.Printf("Liquidity threshold violation for pool %s.\n", poolID)
			automation.handleLiquidityViolation(poolID)
		}
	}
}

// enforceWithdrawalLimits ensures that liquidity pool withdrawals do not exceed set limits
func (automation *DEXEnforcementAutomation) enforceWithdrawalLimits() {
	for _, poolID := range automation.dexManager.GetLiquidityPools() {
		withdrawals := automation.dexManager.GetPoolWithdrawals(poolID)

		if withdrawals > MaxLiquidityWithdrawals {
			fmt.Printf("Withdrawal limit violation for pool %s.\n", poolID)
			automation.handleWithdrawalViolation(poolID)
		}
	}
}

// handleTransactionViolation restricts user activity on exceeding transaction frequency limits
func (automation *DEXEnforcementAutomation) handleTransactionViolation(userID string) {
	err := automation.dexManager.RestrictUserActivity(userID)
	if err != nil {
		fmt.Printf("Failed to restrict user %s due to transaction frequency violation: %v\n", userID, err)
		automation.logDEXAction(userID, "Failed Transaction Restriction")
	} else {
		fmt.Printf("User %s restricted due to excessive transaction frequency.\n", userID)
		automation.logDEXAction(userID, "Transaction Frequency Violation: User Restricted")
	}
}

// handleLiquidityViolation restricts additional withdrawals for pools not meeting liquidity requirements
func (automation *DEXEnforcementAutomation) handleLiquidityViolation(poolID string) {
	err := automation.dexManager.RestrictPoolActivity(poolID)
	if err != nil {
		fmt.Printf("Failed to restrict pool %s due to liquidity threshold violation: %v\n", poolID, err)
		automation.logDEXAction(poolID, "Failed Liquidity Restriction")
	} else {
		fmt.Printf("Pool %s restricted due to insufficient liquidity.\n", poolID)
		automation.logDEXAction(poolID, "Liquidity Violation: Pool Restricted")
	}
}

// handleWithdrawalViolation restricts further withdrawals from pools exceeding withdrawal limits
func (automation *DEXEnforcementAutomation) handleWithdrawalViolation(poolID string) {
	err := automation.dexManager.RestrictWithdrawals(poolID)
	if err != nil {
		fmt.Printf("Failed to restrict withdrawals for pool %s: %v\n", poolID, err)
		automation.logDEXAction(poolID, "Failed Withdrawal Restriction")
	} else {
		fmt.Printf("Further withdrawals restricted for pool %s due to excessive withdrawal activity.\n", poolID)
		automation.logDEXAction(poolID, "Withdrawal Limit Violation: Withdrawals Restricted")
	}
}

// logDEXAction securely logs actions related to DEX enforcement
func (automation *DEXEnforcementAutomation) logDEXAction(entityID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Entity ID: %s", action, entityID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("dex-enforcement-%s-%d", entityID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DEX Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log DEX enforcement action for entity %s: %v\n", entityID, err)
	} else {
		fmt.Println("DEX enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DEXEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualDEXComplianceCheck allows administrators to manually check compliance for a specific pool or user
func (automation *DEXEnforcementAutomation) TriggerManualDEXComplianceCheck(entityID string, entityType string) {
	fmt.Printf("Manually triggering DEX compliance check for entity: %s\n", entityID)

	switch entityType {
	case "user":
		automation.verifyTransactionFrequency()
	case "pool":
		automation.checkLiquidityThresholds()
		automation.enforceWithdrawalLimits()
	default:
		fmt.Printf("Unknown entity type for DEX compliance check: %s\n", entityType)
		automation.logDEXAction(entityID, "Manual Compliance Check Failed")
	}
}
