package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/atomic_swap"
	"synnergy_network_demo/common"
)

const (
	AtomicSwapCheckInterval = 1 * time.Minute // Interval for checking atomic swaps
)

// AtomicSwapEnforcementAutomation enforces the execution of atomic swaps
type AtomicSwapEnforcementAutomation struct {
	swapManager      *atomic_swap.AtomicSwapManager
	consensusEngine  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	enforcementMutex *sync.RWMutex
}

// NewAtomicSwapEnforcementAutomation initializes atomic swap enforcement automation
func NewAtomicSwapEnforcementAutomation(swapManager *atomic_swap.AtomicSwapManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *AtomicSwapEnforcementAutomation {
	return &AtomicSwapEnforcementAutomation{
		swapManager:      swapManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
	}
}

// StartAtomicSwapEnforcement starts the enforcement automation in a continuous loop
func (automation *AtomicSwapEnforcementAutomation) StartAtomicSwapEnforcement() {
	ticker := time.NewTicker(AtomicSwapCheckInterval)

	go func() {
		for range ticker.C {
			automation.enforceAtomicSwaps()
		}
	}()
}

// enforceAtomicSwaps checks for pending atomic swaps and enforces their execution
func (automation *AtomicSwapEnforcementAutomation) enforceAtomicSwaps() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Get all pending atomic swaps
	pendingSwaps := automation.swapManager.GetPendingAtomicSwaps()

	for _, swap := range pendingSwaps {
		if automation.consensusEngine.ValidateTransaction(swap.TransactionID) {
			automation.executeAtomicSwap(swap)
		} else {
			fmt.Printf("Failed to validate transaction for swap ID: %s\n", swap.SwapID)
		}
	}
}

// executeAtomicSwap finalizes an atomic swap and logs it in the ledger
func (automation *AtomicSwapEnforcementAutomation) executeAtomicSwap(swap atomic_swap.AtomicSwap) {
	// Execute the swap through the swap manager
	err := automation.swapManager.ExecuteSwap(swap)
	if err != nil {
		fmt.Printf("Failed to execute atomic swap: %s, Error: %v\n", swap.SwapID, err)
		return
	}

	// Log the successful swap execution in the ledger
	automation.logAtomicSwap(swap)
}

// logAtomicSwap securely logs atomic swap events into the ledger
func (automation *AtomicSwapEnforcementAutomation) logAtomicSwap(swap atomic_swap.AtomicSwap) {
	entryDetails := fmt.Sprintf("Atomic swap %s successfully executed between %s and %s", swap.SwapID, swap.Sender, swap.Receiver)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("atomic-swap-%s-%d", swap.SwapID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Atomic Swap Execution",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log atomic swap %s in the ledger: %v\n", swap.SwapID, err)
	} else {
		fmt.Println("Atomic swap successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AtomicSwapEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualAtomicSwap allows administrators to manually trigger an atomic swap
func (automation *AtomicSwapEnforcementAutomation) TriggerManualAtomicSwap(swapID string) {
	fmt.Printf("Manually triggering atomic swap: %s\n", swapID)

	swap, err := automation.swapManager.GetSwapByID(swapID)
	if err != nil {
		fmt.Printf("Failed to fetch swap with ID: %s, Error: %v\n", swapID, err)
		return
	}

	// Execute the swap and log it
	automation.executeAtomicSwap(swap)
}
