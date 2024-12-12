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
	LiquidityProvisionCheckInterval = 10 * time.Second // Interval for checking liquidity provision compliance
	MaxLiquidityProvisionAttempts   = 5                // Maximum allowed provision attempts per user
	MinimumLiquidity                = 10000            // Minimum liquidity requirement for liquidity pools
)

// LiquidityProvisionRestrictionAutomation monitors and restricts liquidity provision across the network
type LiquidityProvisionRestrictionAutomation struct {
	consensusSystem             *consensus.SynnergyConsensus
	ledgerInstance              *ledger.Ledger
	stateMutex                  *sync.RWMutex
	liquidityProvisionAttempts  map[string]int // Tracks liquidity provision attempts per user
	liquidityPools              map[string]int // Tracks liquidity in different pools
}

// NewLiquidityProvisionRestrictionAutomation initializes and returns an instance of LiquidityProvisionRestrictionAutomation
func NewLiquidityProvisionRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *LiquidityProvisionRestrictionAutomation {
	return &LiquidityProvisionRestrictionAutomation{
		consensusSystem:            consensusSystem,
		ledgerInstance:             ledgerInstance,
		stateMutex:                 stateMutex,
		liquidityProvisionAttempts: make(map[string]int),
		liquidityPools:             make(map[string]int),
	}
}

// StartLiquidityProvisionMonitoring starts continuous monitoring of liquidity provision processes
func (automation *LiquidityProvisionRestrictionAutomation) StartLiquidityProvisionMonitoring() {
	ticker := time.NewTicker(LiquidityProvisionCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorLiquidityProvision()
		}
	}()
}

// monitorLiquidityProvision checks liquidity pools and restricts provision if necessary
func (automation *LiquidityProvisionRestrictionAutomation) monitorLiquidityProvision() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch liquidity pool data from Synnergy Consensus
	poolData := automation.consensusSystem.GetLiquidityPoolData()

	for poolID, liquidity := range poolData {
		// Ensure minimum liquidity requirements are met
		if liquidity < MinimumLiquidity {
			automation.flagLiquidityViolation(poolID, liquidity, "Liquidity below minimum requirement")
		}
	}

	// Check user-specific provision attempts and apply restrictions if necessary
	for userID, attempts := range automation.liquidityProvisionAttempts {
		if attempts > MaxLiquidityProvisionAttempts {
			automation.flagProvisionViolation(userID, attempts, "Exceeded maximum liquidity provision attempts")
		}
	}
}

// flagLiquidityViolation flags a violation in liquidity levels and logs it in the ledger
func (automation *LiquidityProvisionRestrictionAutomation) flagLiquidityViolation(poolID string, liquidity int, reason string) {
	fmt.Printf("Liquidity violation: Pool ID %s, Liquidity: %d, Reason: %s\n", poolID, liquidity, reason)

	// Log the violation in the ledger
	automation.logLiquidityViolation(poolID, liquidity, reason)
}

// logLiquidityViolation logs the flagged liquidity violation into the ledger with full details
func (automation *LiquidityProvisionRestrictionAutomation) logLiquidityViolation(poolID string, liquidity int, violationReason string) {
	// Create a ledger entry for liquidity violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("liquidity-violation-%s-%d", poolID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Liquidity Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Liquidity pool %s violated the provision rule. Current Liquidity: %d. Reason: %s", poolID, liquidity, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptProvisionData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log liquidity violation:", err)
	} else {
		fmt.Println("Liquidity violation logged.")
	}
}

// flagProvisionViolation flags a user's liquidity provision violation and logs it in the ledger
func (automation *LiquidityProvisionRestrictionAutomation) flagProvisionViolation(userID string, attempts int, reason string) {
	fmt.Printf("Provision violation: User ID %s, Attempts: %d, Reason: %s\n", userID, attempts, reason)

	// Log the violation in the ledger
	automation.logProvisionViolation(userID, attempts, reason)
}

// logProvisionViolation logs the flagged provision violation into the ledger with full details
func (automation *LiquidityProvisionRestrictionAutomation) logProvisionViolation(userID string, attempts int, violationReason string) {
	// Create a ledger entry for provision violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("provision-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Provision Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated the provision rule. Attempts: %d. Reason: %s", userID, attempts, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptProvisionData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log provision violation:", err)
	} else {
		fmt.Println("Provision violation logged.")
	}
}

// encryptProvisionData encrypts the provision data before logging for security
func (automation *LiquidityProvisionRestrictionAutomation) encryptProvisionData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting provision data:", err)
		return data
	}
	return string(encryptedData)
}

