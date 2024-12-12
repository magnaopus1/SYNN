package rollups

import (

	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewRollupFeeManager initializes a new RollupFeeManager
func NewRollupFeeManager(feeID string, baseFee float64, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.RollupFeeManager {
	return &common.RollupFeeManager{
		FeeID:          feeID,
		BaseFee:        baseFee,
		TransactionFees: make(map[string]float64),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
	}
}

// CalculateFee calculates the fee for a given transaction based on its size and base fee
func (fm *common.RollupFeeManager) CalculateFee(tx *common.Transaction) (float64, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Example: The fee can be proportional to the transaction size or other metrics
	fee := fm.BaseFee + float64(len(tx.Data))*0.001 // Base fee + size-based fee

	// Log the fee calculation in the ledger
	err := fm.Ledger.RecordFeeCalculation(tx.TxID, fee, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to log fee calculation: %v", err)
	}

	fmt.Printf("Fee for transaction %s calculated as %f\n", tx.TxID, fee)
	return fee, nil
}

// ApplyFee applies the fee to a transaction and stores it in the manager
func (fm *common.RollupFeeManager) ApplyFee(tx *common.Transaction) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Calculate the fee for the transaction
	fee, err := fm.CalculateFee(tx)
	if err != nil {
		return err
	}

	// Store the fee for the transaction
	fm.TransactionFees[tx.TxID] = fee
	fm.TotalFees += fee

	// Encrypt the fee data
	encryptedFee, err := fm.Encryption.EncryptData([]byte(fmt.Sprintf("%f", fee)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt fee data: %v", err)
	}
	tx.Fee = string(encryptedFee)

	// Log the fee application in the ledger
	err = fm.Ledger.RecordFeeApplication(tx.TxID, fee, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log fee application: %v", err)
	}

	fmt.Printf("Fee of %f applied to transaction %s\n", fee, tx.TxID)
	return nil
}

// ValidateFees ensures all fees are correctly applied and accounted for using Synnergy Consensus
func (fm *common.RollupFeeManager) ValidateFees() error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Use consensus to validate all the stored transaction fees
	err := fm.Consensus.ValidateFees(fm.FeeID, fm.TransactionFees)
	if err != nil {
		return fmt.Errorf("fee validation failed: %v", err)
	}

	// Log the fee validation in the ledger
	err = fm.Ledger.RecordFeeValidation(fm.FeeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log fee validation: %v", err)
	}

	fmt.Printf("Fees validated for FeeID %s\n", fm.FeeID)
	return nil
}

// RetrieveFee retrieves the fee for a specific transaction
func (fm *common.RollupFeeManager) RetrieveFee(txID string) (float64, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fee, exists := fm.TransactionFees[txID]
	if !exists {
		return 0, fmt.Errorf("transaction %s has no fee recorded", txID)
	}

	fmt.Printf("Retrieved fee for transaction %s: %f\n", txID, fee)
	return fee, nil
}

// SetBaseFee allows updating the base fee for the rollup
func (fm *common.RollupFeeManager) SetBaseFee(newBaseFee float64) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.BaseFee = newBaseFee

	// Log the base fee update in the ledger
	err := fm.Ledger.RecordBaseFeeUpdate(fm.FeeID, newBaseFee, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log base fee update: %v", err)
	}

	fmt.Printf("Base fee updated to %f for FeeID %s\n", newBaseFee, fm.FeeID)
	return nil
}

// RetrieveBaseFee retrieves the current base fee
func (fm *common.RollupFeeManager) RetrieveBaseFee() float64 {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fmt.Printf("Retrieved base fee: %f for FeeID %s\n", fm.BaseFee, fm.FeeID)
	return fm.BaseFee
}

// RetrieveTotalFees retrieves the total fees collected within the rollup
func (fm *common.RollupFeeManager) RetrieveTotalFees() float64 {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fmt.Printf("Total fees collected for FeeID %s: %f\n", fm.FeeID, fm.TotalFees)
	return fm.TotalFees
}

// RefundFee refunds the fee for a specific transaction
func (fm *common.RollupFeeManager) RefundFee(txID string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fee, exists := fm.TransactionFees[txID]
	if !exists {
		return fmt.Errorf("no fee recorded for transaction %s", txID)
	}

	// Refund the fee
	fm.TotalFees -= fee
	delete(fm.TransactionFees, txID)

	// Log the refund in the ledger
	err := fm.Ledger.RecordFeeRefund(txID, fee, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log fee refund: %v", err)
	}

	fmt.Printf("Refunded fee of %f for transaction %s\n", fee, txID)
	return nil
}
