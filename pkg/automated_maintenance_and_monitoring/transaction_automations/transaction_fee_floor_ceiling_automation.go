package automations

import (
	"fmt"
	"log"
	"sync"
	"time"

	"synnergy_network_demo/transactions"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
)

// FeeCeilingFloorAutomation automates the enforcement of fee floors and ceilings for every transaction in the pool before validation.
type FeeCeilingFloorAutomation struct {
	ledgerInstance     *ledger.Ledger
	transactionPool    *transactions.TransactionPool
	feeManager         *transactions.FeeManager
	mutex              sync.Mutex
	stopChan           chan bool
	adjustmentInterval time.Duration
	networkLoadMonitor func() float64 // Function to monitor network load for dynamic sub-ceiling adjustments
}

// NewFeeCeilingFloorAutomation initializes the FeeCeilingFloorAutomation.
func NewFeeCeilingFloorAutomation(ledgerInstance *ledger.Ledger, transactionPool *transactions.TransactionPool, feeManager *transactions.FeeManager, networkLoadMonitor func() float64) *FeeCeilingFloorAutomation {
	return &FeeCeilingFloorAutomation{
		ledgerInstance:     ledgerInstance,
		transactionPool:    transactionPool,
		feeManager:         feeManager,
		adjustmentInterval: 50 * time.Millisecond, // Run every 50 milliseconds
		networkLoadMonitor: networkLoadMonitor,
		stopChan:           make(chan bool),
	}
}

// Start begins the continuous process of checking fees and enforcing ceilings/floors for each transaction in the pool before validation.
func (f *FeeCeilingFloorAutomation) Start() {
	go f.runFeeEnforcementLoop()
	log.Println("Transaction Fee Ceiling/Floor Automation started.")
}

// Stop stops the continuous process of fee enforcement.
func (f *FeeCeilingFloorAutomation) Stop() {
	f.stopChan <- true
	log.Println("Transaction Fee Ceiling/Floor Automation stopped.")
}

// runFeeEnforcementLoop periodically checks and enforces fee floors and ceilings for each transaction in the pool before validation.
func (f *FeeCeilingFloorAutomation) runFeeEnforcementLoop() {
	ticker := time.NewTicker(f.adjustmentInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f.enforceFeesForAllTransactions()
		case <-f.stopChan:
			return
		}
	}
}

// enforceFeesForAllTransactions checks and enforces fee floors and ceilings for every transaction in the transaction pool.
func (f *FeeCeilingFloorAutomation) enforceFeesForAllTransactions() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Fetch all transactions from the pool
	transactions, err := f.transactionPool.ListTransactions()
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return
	}

	// Enforce the fee ceiling and floor for each transaction
	for _, tx := range transactions {
		err := f.EnforceFeeForTransaction(&tx)
		if err != nil {
			log.Printf("Failed to enforce fee for transaction %s: %v", tx.ID, err)
		}
	}
}

// EnforceFeeForTransaction enforces the fixed ceiling, dynamic sub-ceiling, and fee floor for a specific transaction.
func (f *FeeCeilingFloorAutomation) EnforceFeeForTransaction(transaction *common.Transaction) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fee := transaction.Fee
	amount := transaction.Amount

	// Retrieve current floor and sub-ceiling values
	feeFloor, err := f.feeManager.GetFeeFloor()
	if err != nil {
		return fmt.Errorf("failed to retrieve fee floor: %v", err)
	}

	feeCeiling, err := f.feeManager.GetFeeCeiling()
	if err != nil {
		return fmt.Errorf("failed to retrieve fee sub-ceiling: %v", err)
	}

	// Enforce fee floor
	if fee < feeFloor {
		fee = feeFloor
	}

	// Enforce dynamic sub-ceiling
	if fee > feeCeiling {
		fee = feeCeiling
	}

	// Enforce the fixed 0.25% ceiling
	if fee > amount*transactions.FeeCeilingPercent {
		fee = amount * transactions.FeeCeilingPercent
	}

	transaction.Fee = fee
	log.Printf("Transaction %s fee enforced. Final fee: %f (Floor: %f, Sub-Ceiling: %f, Fixed Ceiling: %f)", transaction.ID, fee, feeFloor, feeCeiling, transactions.FeeCeilingPercent)

	// Encrypt transaction fee before adding to ledger
	encryptedFee, err := encryption.EncryptData(fmt.Sprintf("%f", fee), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction fee: %v", err)
	}
	transaction.FeeEncrypted = encryptedFee

	// Add the transaction to the ledger with the enforced fee
	err = f.ledgerInstance.AddTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to add transaction to ledger: %v", err)
	}

	return nil
}

// adjustFeeFloorsAndCeilings adjusts the fee floors and dynamic sub-ceilings based on network load.
func (f *FeeCeilingFloorAutomation) adjustFeeFloorsAndCeilings() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	networkLoad := f.networkLoadMonitor()

	// Adjust the dynamic sub-ceiling based on network load, but ensure it doesn't exceed the fixed ceiling of 0.25%.
	feeSubCeiling := f.calculateSubCeiling(networkLoad)
	feeFloor := f.calculateFeeFloor(networkLoad)

	// Apply the new floor and sub-ceiling to the fee manager
	err := f.feeManager.SetFeeFloor(feeFloor)
	if err != nil {
		log.Printf("Failed to set fee floor: %v", err)
	}

	err = f.feeManager.SetFeeCeiling(feeSubCeiling) // Setting the dynamic sub-ceiling
	if err != nil {
		log.Printf("Failed to set fee sub-ceiling: %v", err)
	}

	log.Printf("Fee floor set to %f and dynamic sub-ceiling set to %f based on network load: %f", feeFloor, feeSubCeiling, networkLoad)
}

// calculateSubCeiling dynamically adjusts the sub-ceiling but ensures it does not exceed the fixed ceiling of 0.25%.
func (f *FeeCeilingFloorAutomation) calculateSubCeiling(networkLoad float64) float64 {
	// Base dynamic sub-ceiling logic: range between FeeFloorPercent and FeeCeilingPercent, adjusting based on load.
	baseSubCeiling := transactions.FeeCeilingPercent * (1 - (networkLoad / 2)) // Less ceiling under high load.
	return f.min(transactions.FeeCeilingPercent, baseSubCeiling)               // Never exceed 0.25%
}

// calculateFeeFloor dynamically adjusts the fee floor based on network load.
func (f *FeeCeilingFloorAutomation) calculateFeeFloor(networkLoad float64) float64 {
	// Basic example of fee floor logic. Increase floor under high load to discourage micro-transactions.
	baseFloor := transactions.FeeFloorPercent
	return baseFloor * (1 + networkLoad) // Increase floor with load
}

// Utility functions to ensure values stay within bounds
func (f *FeeCeilingFloorAutomation) min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
