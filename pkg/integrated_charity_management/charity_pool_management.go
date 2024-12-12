package integrated_charity_management

import (
	"encoding/base64"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewCharityPoolManagement initializes the charity pool manager
func NewCharityPoolManagement(ledgerInstance *ledger.Ledger) *CharityPoolManagement {
    return &CharityPoolManagement{
        InternalPoolBalance: 0, // Initial balance is 0
        ExternalPoolBalance: 0, // Initial balance is 0
        LedgerInstance:      ledgerInstance,
    }
}

// UpdateCharityPools receives a portion of the transaction fees and apportions them equally
func (cpm *CharityPoolManagement) UpdateCharityPools(transactionFee float64) error {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    if transactionFee <= 0 {
        return fmt.Errorf("invalid transaction fee amount: %.2f", transactionFee)
    }

    // Split the transaction fee equally between the two pools
    halfFee := transactionFee / 2
    cpm.InternalPoolBalance += halfFee
    cpm.ExternalPoolBalance += halfFee

    fmt.Printf("Transaction fee of %.2f SYNN apportioned: %.2f to internal pool, %.2f to external pool.\n",
        transactionFee, halfFee, halfFee)

    // Log the fee distribution to the ledger
    err := cpm.logCharityFeeDistributionToLedger(transactionFee, halfFee)
    if err != nil {
        return fmt.Errorf("failed to log charity fee distribution to ledger: %v", err)
    }

    return nil
}

// GetInternalPoolBalance returns the current balance of the internal charity pool
func (cpm *CharityPoolManagement) GetInternalPoolBalance() float64 {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    return cpm.InternalPoolBalance
}

// GetExternalPoolBalance returns the current balance of the external charity pool
func (cpm *CharityPoolManagement) GetExternalPoolBalance() float64 {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    return cpm.ExternalPoolBalance
}

// WithdrawFromInternalPool handles withdrawals from the internal charity pool
func (cpm *CharityPoolManagement) WithdrawFromInternalPool(amount float64) error {
    cpm.mutex.Lock()
    defer cpm.mutex.Unlock()

    if amount > cpm.InternalPoolBalance {
        return fmt.Errorf("insufficient balance in internal charity pool")
    }

    cpm.InternalPoolBalance -= amount
    fmt.Printf("Withdrawn %.2f SYNN from internal charity pool. Remaining balance: %.2f SYNN.\n", amount, cpm.InternalPoolBalance)

    // Log the withdrawal to the ledger
    return cpm.logWithdrawalToLedger("Internal Charity Pool", amount)
}


// logCharityFeeDistributionToLedger logs the distribution of fees between the internal and external pools
func (cpm *CharityPoolManagement) logCharityFeeDistributionToLedger(transactionFee, halfFee float64) error {
    logData := fmt.Sprintf("Transaction fee: %.2f SYNN, Internal Pool: %.2f SYNN, External Pool: %.2f SYNN",
        transactionFee, halfFee, halfFee)

    // Step 1: Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming NewEncryption creates AES with 256-bit key
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Step 2: Encrypt the log data using the encryption instance
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logData), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt charity pool distribution log: %v", err)
    }

    // Step 3: Convert the encrypted log to a base64-encoded string (if you need to log it separately)
    encryptedLogString := base64.StdEncoding.EncodeToString(encryptedLog)

    // Step 4: Record the transaction fee to the ledger (removing the encrypted log string argument)
    cpm.LedgerInstance.RecordCharityFeeDistribution(transactionFee)

    // Optionally, you can log the encrypted data elsewhere if needed
    fmt.Printf("Encrypted Charity Fee Distribution Log: %s\n", encryptedLogString)

    return nil
}

// logWithdrawalToLedger logs a withdrawal from either the internal or external charity pool
func (cpm *CharityPoolManagement) logWithdrawalToLedger(poolName string, amount float64) error {
    logData := fmt.Sprintf("%s: Withdrawal of %.2f SYNN", poolName, amount)

    // Step 1: Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Assuming NewEncryption creates AES with a 256-bit key
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Step 2: Encrypt the log data using the encryption instance
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logData), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt charity pool withdrawal log: %v", err)
    }

    // Step 3: Convert the encrypted log to a base64-encoded string (if you need to log it separately)
    encryptedLogString := base64.StdEncoding.EncodeToString(encryptedLog)

    // Step 4: Record the withdrawal amount to the ledger (removing the encrypted log string argument)
    cpm.LedgerInstance.RecordCharityPoolWithdrawal(amount)

    // Optionally, you can log the encrypted data elsewhere if needed
    fmt.Printf("Encrypted Charity Pool Withdrawal Log: %s\n", encryptedLogString)

    return nil
}
