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

// Constants for transaction volume restrictions
const (
	TransactionVolumeCheckInterval = 2 * time.Second  // Interval for checking transaction volume
	MaxSubBlockTransactionVolume   = 1000             // Maximum transaction volume per sub-block
	MaxBlockTransactionVolume      = 1000000           // Maximum transaction volume per block
)

// TransactionVolumeRestrictionAutomation monitors and restricts transaction volume
type TransactionVolumeRestrictionAutomation struct {
	consensusSystem  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	transactionMutex *sync.RWMutex
}

// NewTransactionVolumeRestrictionAutomation initializes the automation for monitoring transaction volume
func NewTransactionVolumeRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *TransactionVolumeRestrictionAutomation {
	return &TransactionVolumeRestrictionAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		transactionMutex: transactionMutex,
	}
}

// StartMonitoring begins continuous monitoring of transaction volumes in sub-blocks and blocks
func (automation *TransactionVolumeRestrictionAutomation) StartMonitoring() {
	ticker := time.NewTicker(TransactionVolumeCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateTransactionVolumes()
		}
	}()
}

// evaluateTransactionVolumes checks if any sub-blocks or blocks exceed the maximum transaction volume
func (automation *TransactionVolumeRestrictionAutomation) evaluateTransactionVolumes() {
	automation.transactionMutex.Lock()
	defer automation.transactionMutex.Unlock()

	// Check sub-block transaction volume
	subBlocks := automation.consensusSystem.GetPendingSubBlocks()
	for _, subBlock := range subBlocks {
		if subBlock.Volume > MaxSubBlockTransactionVolume {
			automation.logVolumeViolation(subBlock.ID, subBlock.Volume, "Sub-block")
			automation.applyVolumeReduction(subBlock.ID, MaxSubBlockTransactionVolume)
		}
	}

	// Check block transaction volume
	blocks := automation.consensusSystem.GetPendingBlocks()
	for _, block := range blocks {
		if block.Volume > MaxBlockTransactionVolume {
			automation.logVolumeViolation(block.ID, block.Volume, "Block")
			automation.applyVolumeReduction(block.ID, MaxBlockTransactionVolume)
		}
	}
}

// logVolumeViolation logs blocks or sub-blocks that exceed the volume limits in the ledger
func (automation *TransactionVolumeRestrictionAutomation) logVolumeViolation(id string, volume int, entityType string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("volume-violation-%s-%d", id, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      fmt.Sprintf("%s Transaction Volume Violation", entityType),
		Status:    "Violation",
		Details:   fmt.Sprintf("%s %s exceeded the volume limit with %d transactions.", entityType, id, volume),
	}

	// Encrypt the violation details before storing in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log volume violation:", err)
	} else {
		fmt.Println("Transaction volume violation logged for:", id)
	}
}

// applyVolumeReduction applies reductions to transaction volumes that exceed the ceiling
func (automation *TransactionVolumeRestrictionAutomation) applyVolumeReduction(id string, limit int) {
	fmt.Printf("Applying volume reduction for %s to adhere to the transaction volume limit.\n", id)

	// Adjust the volume to meet the limit in the consensus system
	err := automation.consensusSystem.AdjustTransactionVolume(id, limit)
	if err != nil {
		fmt.Println("Failed to adjust transaction volume for:", id, "Error:", err)
		return
	}

	// Log the volume adjustment in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("volume-adjustment-%s-%d", id, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Volume Adjustment",
		Status:    "Adjusted",
		Details:   fmt.Sprintf("Volume for %s reduced to %d transactions due to exceeding the limit.", id, limit),
	}

	// Encrypt the adjustment details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log volume adjustment:", err)
	} else {
		fmt.Println("Volume adjustment logged for:", id)
	}
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *TransactionVolumeRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
