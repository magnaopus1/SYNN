package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
)

// Configuration for transaction fee enforcement automation
const (
	FeeCheckInterval           = 10 * time.Second // Interval to check and enforce transaction fee policies
	MinTransactionFee          = 0.001            // Minimum transaction fee
	MaxTransactionFee          = 0.01             // Maximum transaction fee to control congestion
	CongestionThresholdTPS     = 4000             // TPS level at which fees should increase
	CongestionFeeMultiplier    = 1.5              // Multiplier to increase fee during high congestion
)

// TransactionFeeEnforcementAutomation monitors and enforces transaction fee standards
type TransactionFeeEnforcementAutomation struct {
	networkManager        *network.NetworkManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	currentTransactionFee float64 // Current enforced transaction fee
}

// NewTransactionFeeEnforcementAutomation initializes the transaction fee enforcement automation
func NewTransactionFeeEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *TransactionFeeEnforcementAutomation {
	return &TransactionFeeEnforcementAutomation{
		networkManager:        networkManager,
		consensusEngine:       consensusEngine,
		ledgerInstance:        ledgerInstance,
		enforcementMutex:      enforcementMutex,
		currentTransactionFee: MinTransactionFee,
	}
}

// StartTransactionFeeEnforcement begins continuous monitoring and adjustment of transaction fees
func (automation *TransactionFeeEnforcementAutomation) StartTransactionFeeEnforcement() {
	ticker := time.NewTicker(FeeCheckInterval)

	go func() {
		for range ticker.C {
			automation.adjustTransactionFee()
		}
	}()
}

// adjustTransactionFee monitors the TPS and adjusts the transaction fee as necessary
func (automation *TransactionFeeEnforcementAutomation) adjustTransactionFee() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	currentTPS := automation.networkManager.GetCurrentTPS()
	fmt.Printf("Current TPS: %d, Current Fee: %f\n", currentTPS, automation.currentTransactionFee)

	if currentTPS > CongestionThresholdTPS {
		automation.increaseTransactionFee()
	} else {
		automation.resetTransactionFee()
	}
}

// increaseTransactionFee increases the transaction fee during high congestion
func (automation *TransactionFeeEnforcementAutomation) increaseTransactionFee() {
	newFee := automation.currentTransactionFee * CongestionFeeMultiplier
	if newFee > MaxTransactionFee {
		newFee = MaxTransactionFee
	}

	if newFee != automation.currentTransactionFee {
		fmt.Printf("Increasing transaction fee from %f to %f due to high congestion.\n", automation.currentTransactionFee, newFee)
		automation.currentTransactionFee = newFee
		automation.logFeeAdjustment("Fee Increased", fmt.Sprintf("New Fee: %f due to congestion", newFee))
	}
}

// resetTransactionFee resets the transaction fee to the minimum when congestion is low
func (automation *TransactionFeeEnforcementAutomation) resetTransactionFee() {
	if automation.currentTransactionFee != MinTransactionFee {
		fmt.Printf("Resetting transaction fee to minimum: %f\n", MinTransactionFee)
		automation.currentTransactionFee = MinTransactionFee
		automation.logFeeAdjustment("Fee Reset", "Reset to minimum due to normal traffic")
	}
}

// logFeeAdjustment securely logs actions related to transaction fee adjustments
func (automation *TransactionFeeEnforcementAutomation) logFeeAdjustment(action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Details: %s", action, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("transaction-fee-enforcement-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Fee Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log transaction fee enforcement action: %v\n", err)
	} else {
		fmt.Println("Transaction fee enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *TransactionFeeEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualFeeAdjustment allows administrators to manually set a transaction fee
func (automation *TransactionFeeEnforcementAutomation) TriggerManualFeeAdjustment(newFee float64) {
	if newFee < MinTransactionFee || newFee > MaxTransactionFee {
		fmt.Printf("Invalid fee adjustment: %f. Must be within %f and %f.\n", newFee, MinTransactionFee, MaxTransactionFee)
		automation.logFeeAdjustment("Manual Fee Adjustment Failed", "Fee outside allowed range")
		return
	}

	fmt.Printf("Manually adjusting transaction fee to: %f\n", newFee)
	automation.currentTransactionFee = newFee
	automation.logFeeAdjustment("Manual Fee Adjusted", fmt.Sprintf("New Fee: %f", newFee))
}
