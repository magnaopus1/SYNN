package transactions

import (
	"fmt"
	"math"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// FeeManager manages the calculation and enforcement of transaction fee ceilings and floors.
type FeeManager struct {
	ledgerInstance       *ledger.Ledger // Reference to the ledger instance
	mutex                sync.Mutex      // Mutex for thread-safe operations
	BaseFee              float64         // The base fee for transactions
	NetworkLoad          float64         // Current network load percentage (0 to 1)
	PendingTransactions   int             // Current number of pending transactions
}

// Constants defining the ceiling and floor of transaction fees.
const (
	FeeCeilingPercent = 0.005     // 0.5% of the transaction amount
	FeeFloorPercent   = 0.000000001 // 0.0000001% of the transaction amount
)

// NewFeeManager initializes a new FeeManager with the given ledger instance.
func NewFeeManager(ledgerInstance *ledger.Ledger) *FeeManager {
	return &FeeManager{
		ledgerInstance: ledgerInstance,
		BaseFee:       0.0001, // Default base fee, adjust as needed
	}
}

// CalculateTransactionFee calculates the fee for a given transaction amount, enforcing ceiling and floor limits.
func (fm *FeeManager) CalculateTransactionFee(amount float64) (float64, error) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	if amount <= 0 {
		return 0, fmt.Errorf("transaction amount must be greater than zero")
	}

	// Calculate fee based on percentage.
	fee := amount * FeeCeilingPercent

	// Ensure the fee is not below the floor.
	if fee < amount*FeeFloorPercent {
		fee = amount * FeeFloorPercent
	}

	// Ensure the fee is not above the ceiling.
	if fee > amount*FeeCeilingPercent {
		fee = amount * FeeCeilingPercent
	}

	// Log the calculated fee in the ledger for auditability.
	err := fm.logTransactionFee(amount, fee)
	if err != nil {
		return 0, fmt.Errorf("failed to log transaction fee: %v", err)
	}

	return fee, nil
}

// logTransactionFee logs the calculated fee to the ledger, encrypting the details for security.
func (fm *FeeManager) logTransactionFee(amount float64, fee float64) error {
	// Prepare the log entry.
	logEntry := fmt.Sprintf("Transaction amount: %f, Fee: %f", amount, fee)

	// Create an encryption instance, handle two return values.
	encryptionInstance, err := common.NewEncryption(256) // Assuming 256 bits, adjust if needed.
	if err != nil {
		return fmt.Errorf("failed to initialize encryption instance: %v", err)
	}

	// Encrypt the log entry.
	encryptedLogEntry, err := encryptionInstance.EncryptData("AES", []byte(logEntry), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction fee log: %v", err)
	}

	// Convert the encrypted log entry to a string.
	encryptedLogStr := string(encryptedLogEntry)

	// Record the log entry in the ledger. Assuming the second parameter is a timestamp (uint64).
	timestamp := uint64(time.Now().Unix())
	err = fm.ledgerInstance.RecordFeeLog(encryptedLogStr, timestamp)
	if err != nil {
		return fmt.Errorf("failed to record transaction fee log in ledger: %v", err)
	}

	return nil
}

// ValidateFee checks if a calculated fee is within acceptable boundaries and returns an error if it violates any rules.
func (fm *FeeManager) ValidateFee(amount float64, fee float64) error {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	// Enforce fee ceiling and floor constraints.
	if fee < amount*FeeFloorPercent || fee > amount*FeeCeilingPercent {
		return fmt.Errorf("fee %f for transaction amount %f violates ceiling/floor rules", fee, amount)
	}

	return nil
}

// EnforceFeeCeiling enforces the maximum fee ceiling on a transaction, ensuring it does not exceed the allowed maximum.
func (fm *FeeManager) EnforceFeeCeiling(amount float64, fee float64) (float64, error) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	// If the fee exceeds the ceiling, adjust it to the ceiling.
	if fee > amount*FeeCeilingPercent {
		fee = amount * FeeCeilingPercent
	}

	return fee, nil
}

// EnforceFeeFloor ensures the fee does not fall below the minimum acceptable fee.
func (fm *FeeManager) EnforceFeeFloor(amount float64, fee float64) (float64, error) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	// If the fee is below the floor, adjust it to the floor.
	if fee < amount*FeeFloorPercent {
		fee = amount * FeeFloorPercent
	}

	return fee, nil
}

// AdjustFee adjusts the fee for a transaction based on blockchain conditions like network load or validator requirements.
func (fm *FeeManager) AdjustFee(amount float64) (float64, error) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	// Calculate the base fee
	baseFee := fm.BaseFee

	// Adjust base fee based on network load
	if fm.NetworkLoad > 0.75 { // If the network load is high
		baseFee *= 1.5 // Increase fee by 50%
	} else if fm.NetworkLoad < 0.25 { // If the network load is low
		baseFee *= 0.75 // Decrease fee by 25%
	}

	// Adjust base fee based on the number of pending transactions
	if fm.PendingTransactions > 100 { // Arbitrary threshold for pending transactions
		baseFee *= 1.2 // Increase fee by 20%
	} else if fm.PendingTransactions < 20 {
		baseFee *= 0.8 // Decrease fee by 20%
	}

	// Ensure adjusted fee is within the defined floor and ceiling
	adjustedFee := math.Max(baseFee, amount*FeeFloorPercent)
	adjustedFee = math.Min(adjustedFee, amount*FeeCeilingPercent)

	return adjustedFee, nil
}
