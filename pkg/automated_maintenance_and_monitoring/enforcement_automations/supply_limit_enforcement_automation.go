package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/token"
)

// Configuration for supply limit enforcement automation
const (
	SupplyCheckInterval         = 15 * time.Second // Interval to check for supply limit violations
	MaxSupplyViolationThreshold = 1                // Maximum allowed violations before enforcement action
)

// SupplyLimitEnforcementAutomation monitors and enforces supply limits for tokens and assets
type SupplyLimitEnforcementAutomation struct {
	tokenManager         *token.TokenManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	supplyViolationMap   map[string]int // Tracks supply limit violation count for each token or asset
}

// NewSupplyLimitEnforcementAutomation initializes the supply limit enforcement automation
func NewSupplyLimitEnforcementAutomation(tokenManager *token.TokenManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *SupplyLimitEnforcementAutomation {
	return &SupplyLimitEnforcementAutomation{
		tokenManager:       tokenManager,
		consensusEngine:    consensusEngine,
		ledgerInstance:     ledgerInstance,
		enforcementMutex:   enforcementMutex,
		supplyViolationMap: make(map[string]int),
	}
}

// StartSupplyLimitEnforcement begins continuous monitoring and enforcement of supply limits
func (automation *SupplyLimitEnforcementAutomation) StartSupplyLimitEnforcement() {
	ticker := time.NewTicker(SupplyCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkSupplyLimits()
		}
	}()
}

// checkSupplyLimits monitors each token’s supply and enforces actions if necessary
func (automation *SupplyLimitEnforcementAutomation) checkSupplyLimits() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateSupplyLimits()
	automation.enforceSupplyLimits()
}

// evaluateSupplyLimits checks each token for supply violations and flags tokens exceeding limits
func (automation *SupplyLimitEnforcementAutomation) evaluateSupplyLimits() {
	for _, tokenID := range automation.tokenManager.GetAllTokens() {
		supplyExceeded := automation.tokenManager.IsSupplyExceeded(tokenID)
		if supplyExceeded {
			automation.supplyViolationMap[tokenID]++
			fmt.Printf("Token %s has exceeded its supply limit.\n", tokenID)
		} else {
			automation.supplyViolationMap[tokenID] = 0 // Reset if within limits
		}
	}
}

// enforceSupplyLimits takes action on tokens that exceed allowed supply limit violations
func (automation *SupplyLimitEnforcementAutomation) enforceSupplyLimits() {
	for tokenID, violations := range automation.supplyViolationMap {
		if violations > MaxSupplyViolationThreshold {
			fmt.Printf("Enforcing supply limit action on token %s due to supply cap violation.\n", tokenID)
			automation.freezeToken(tokenID)
		}
	}
}

// freezeToken freezes a token that has exceeded its supply limit violations
func (automation *SupplyLimitEnforcementAutomation) freezeToken(tokenID string) {
	err := automation.tokenManager.FreezeToken(tokenID)
	if err != nil {
		fmt.Printf("Failed to freeze token %s: %v\n", tokenID, err)
		automation.logSupplyAction(tokenID, "Freeze Failed", "Exceeded Supply Limit")
	} else {
		fmt.Printf("Token %s has been frozen due to exceeding the supply limit.\n", tokenID)
		automation.logSupplyAction(tokenID, "Token Frozen", "Exceeded Supply Limit")
	}
}

// logSupplyAction securely logs actions related to supply limit enforcement
func (automation *SupplyLimitEnforcementAutomation) logSupplyAction(tokenID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Token ID: %s, Details: %s", action, tokenID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("supply-limit-enforcement-%s-%d", tokenID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Supply Limit Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log supply limit enforcement action for token %s: %v\n", tokenID, err)
	} else {
		fmt.Println("Supply limit enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *SupplyLimitEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualSupplyFreeze allows administrators to manually freeze a token if needed
func (automation *SupplyLimitEnforcementAutomation) TriggerManualSupplyFreeze(tokenID string) {
	fmt.Printf("Manually freezing token due to supply cap: %s\n", tokenID)

	if automation.tokenManager.IsSupplyExceeded(tokenID) {
		automation.freezeToken(tokenID)
	} else {
		fmt.Printf("Token %s is within supply limits, freeze not required.\n", tokenID)
		automation.logSupplyAction(tokenID, "Manual Freeze Skipped", "Within Supply Limits")
	}
}
