package common

import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// GasManager manages gas fees for smart contract executions in the virtual machine.
type GasManager struct {
	LedgerInstance  *ledger.Ledger             // Ledger instance to track gas usage and fees
	mutex           sync.Mutex                 // Mutex for thread-safe operations
	GasPrice        float64                    // Current gas price (in terms of native currency)
    ConsensusEngine *SynnergyConsensus
}

// NewGasManager initializes a new GasManager.
func NewGasManager(ledgerInstance *ledger.Ledger, consensusEngine *SynnergyConsensus, gasPrice float64) *GasManager {
    return &GasManager{
        LedgerInstance:  ledgerInstance,
        ConsensusEngine: consensusEngine,
        GasPrice:        gasPrice,
    }
}

// CalculateGas calculates the gas fee for a given smart contract execution.
func (gm *GasManager) CalculateGas(contractID string, bytecode string, executionTime time.Duration) (float64, error) {
    gm.mutex.Lock()
    defer gm.mutex.Unlock()

    gasUnits := gm.estimateGasUnits(bytecode, executionTime)
    totalGas := gasUnits * gm.GasPrice

    fmt.Printf("Gas calculated for contract %s: %.6f units at price %.6f (total: %.6f)\n", contractID, gasUnits, gm.GasPrice, totalGas)

    return totalGas, nil
}

// estimateGasUnits estimates the gas units required based on bytecode complexity and execution time.
func (gm *GasManager) estimateGasUnits(bytecode string, executionTime time.Duration) float64 {
    bytecodeLength := len(bytecode)
    baseGasUnits := float64(bytecodeLength) * 0.1  // Estimate based on bytecode length
    timeFactor := float64(executionTime.Seconds()) * 10 // Factor based on execution time

    return baseGasUnits + timeFactor
}

// DeductGas deducts the gas fee from the user's balance and records it in the ledger.
func (gm *GasManager) DeductGas(walletAddress string, gasAmount float64) error {
    gm.mutex.Lock()
    defer gm.mutex.Unlock()

    // Check if the wallet has enough balance
    balance, err := gm.LedgerInstance.GetBalance(walletAddress)
    if err != nil {
        return fmt.Errorf("failed to retrieve balance for wallet %s: %v", walletAddress, err)
    }

    if balance < gasAmount {
        return fmt.Errorf("insufficient balance for wallet %s to pay gas (required: %.6f, available: %.6f)", walletAddress, gasAmount, balance)
    }

    // Deduct the gas fee from the wallet
    newBalance := uint64(balance - gasAmount)  // Convert float64 to uint64 if necessary
    err = gm.LedgerInstance.AccountsWalletLedger.UpdateBalance(walletAddress, newBalance)
    if err != nil {
        return fmt.Errorf("failed to deduct gas from wallet %s: %v", walletAddress, err)
    }

    // Record the gas deduction transaction in the ledger (only passing three arguments)
    err = gm.LedgerInstance.BlockchainConsensusCoinLedger.RecordTransaction(walletAddress, "GasPool", gasAmount)
    if err != nil {
        return fmt.Errorf("failed to record gas deduction in the ledger: %v", err)
    }

    fmt.Printf("Gas fee of %.6f deducted from wallet %s, new balance: %.6f\n", gasAmount, walletAddress, balance-gasAmount)
    return nil
}



// RefundGas refunds unused gas to the user's balance after execution is complete.
func (gm *GasManager) RefundGas(walletAddress string, unusedGas float64) error {
    gm.mutex.Lock()
    defer gm.mutex.Unlock()

    // Add the refunded gas back to the wallet
    balance, err := gm.LedgerInstance.GetBalance(walletAddress)
    if err != nil {
        return fmt.Errorf("failed to retrieve balance for wallet %s: %v", walletAddress, err)
    }

    // Update the balance and convert it to uint64 if necessary
    newBalance := balance + unusedGas
    err = gm.LedgerInstance.AccountsWalletLedger.UpdateBalance(walletAddress, uint64(newBalance))  // Convert float64 to uint64 if needed
    if err != nil {
        return fmt.Errorf("failed to refund gas to wallet %s: %v", walletAddress, err)
    }

    // Record the gas refund transaction in the ledger (removed the fourth argument "GAS_REFUND")
    err = gm.LedgerInstance.BlockchainConsensusCoinLedger.RecordTransaction("GasPool", walletAddress, unusedGas)
    if err != nil {
        return fmt.Errorf("failed to record gas refund in the ledger: %v", err)
    }

    fmt.Printf("Unused gas of %.6f refunded to wallet %s, new balance: %.6f\n", unusedGas, walletAddress, newBalance)
    return nil
}



// UpdateGasPrice dynamically updates the gas price based on network conditions or policies.
func (gm *GasManager) UpdateGasPrice(newGasPrice float64) {
    gm.mutex.Lock()
    defer gm.mutex.Unlock()

    gm.GasPrice = newGasPrice
    fmt.Printf("Gas price updated to %.6f\n", newGasPrice)
}

// ChargeGas calculates and deducts gas from the user and records it in the ledger.
func (gm *GasManager) ChargeGas(walletAddress string, contractID string, bytecode string, executionTime time.Duration) error {
    // Step 1: Calculate the gas fee
    gasAmount, err := gm.CalculateGas(contractID, bytecode, executionTime)
    if err != nil {
        return fmt.Errorf("failed to calculate gas: %v", err)
    }

    // Step 2: Deduct the gas from the user's balance
    err = gm.DeductGas(walletAddress, gasAmount)
    if err != nil {
        return fmt.Errorf("failed to deduct gas: %v", err)
    }

    return nil
}

// logGasTransaction logs a gas-related transaction into the ledger.
func (gm *GasManager) logGasTransaction(transactionType string, walletAddress string, gasAmount float64, encryptionInstance *Encryption, encryptionKey []byte) error {
    // Create the transaction details as a string
    transaction := fmt.Sprintf("Wallet: %s, To: GasPool, Amount: %.6f, Type: %s", walletAddress, gasAmount, transactionType)

    // Encrypt the transaction data (if needed)
    _, err := encryptionInstance.EncryptData("AES", []byte(transaction), encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt gas transaction: %v", err)
    }

    // Log the transaction in the ledger (removed the fourth argument "transactionType")
    err = gm.LedgerInstance.BlockchainConsensusCoinLedger.RecordTransaction(walletAddress, "GasPool", gasAmount)
    if err != nil {
        return fmt.Errorf("failed to record gas transaction: %v", err)
    }

    fmt.Printf("Gas transaction logged for wallet %s: %.6f to GasPool\n", walletAddress, gasAmount)
    return nil
}


